package service

import (
	"context"
	"flag"
	"fmt"
	orgsvc "gomicroservices/internal/organization/service"
	"gomicroservices/internal/user/model"
	"gomicroservices/internal/user/repo"
	"gomicroservices/internal/util"
	"time"

	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/pkg/errors"

	"github.com/dgrijalva/jwt-go"
)

var ErrInvalidCredentials = errors.New("invalid credentials")
var ErrInvalidRequest = errors.New("invalid request")

type Service struct {
	repo   repo.DBRepo
	orgsvc orgsvc.IService
}

var repoType string
var sugar string

func init() {
	flag.StringVar(&repoType, "orgsvc", "db", "")
	flag.StringVar(&sugar, "sugar", "", "")
	flag.Parse()
}

func New(db *pgxpool.Pool) Service {
	var o orgsvc.IService
	switch repoType {
	default:
		o = orgsvc.New(db)
	}
	return Service{
		repo:   repo.New(db),
		orgsvc: o,
	}
}

func (s *Service) Login(ctx context.Context, username string, password string) (string, error) {
	user, err := s.repo.GetUserByUsernamePassword(ctx, username, password)
	if err != nil {
		if err == repo.ErrNoUserFound {
			return "", ErrInvalidCredentials
		}
		return "", err

	}
	return generateUserToken(user, sugar)
}

func (s *Service) GetUser(ctx context.Context, id uint) (*model.User, error) {
	log := util.GetLoggerFromContext(ctx)

	log.Infof("Reading user from repo. id=%v", id)
	u, err := s.repo.GetUser(ctx, id)
	if err != nil {
		return nil, err
	}

	log.Infof("Reading branch from org service. id=%s", u.BranchID)
	u.Branch, err = s.orgsvc.GetBranch(ctx, u.BranchID)
	if err != nil {
		return nil, errors.Wrapf(err, "Failed to get branch of user. userId=%v branchId=%v", u.ID, u.BranchID)
	}
	return u, nil
}

func (s *Service) GetUsers(ctx context.Context) (users []*model.User, err error) {
	authUser := ctx.Value("user").(*model.User)

	if authUser.Role == "superadmin" {
		users, err = s.repo.GetUsers(ctx)
	} else {
		users, err = s.repo.GetUsersByOrganization(ctx, authUser.OrganizationID)
	}

	return
}

func (s *Service) CreateUser(ctx context.Context, user *model.User) error {

	if err := s.repo.CreateUser(ctx, user); err != nil {
		if err == repo.ErrFKViolation {
			return ErrInvalidRequest
		}
		return fmt.Errorf("Failed to create user. %v", err)
	}
	return nil
}

func generateUserToken(u *model.User, sugar string) (string, error) {
	claims := jwt.MapClaims{}
	claims["id"] = u.ID
	claims["name"] = u.Name
	claims["username"] = u.Username
	claims["role"] = u.Role
	claims["branchId"] = u.BranchID
	claims["organizationId"] = u.OrganizationID
	claims["exp"] = time.Now().Add(time.Hour * 24).Unix()
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return token.SignedString([]byte(sugar))
}

func AuthCheck(tokenStr string) (*model.User, error) {
	claims := jwt.MapClaims{}
	token, err := jwt.ParseWithClaims(tokenStr, claims, func(token *jwt.Token) (interface{}, error) {
		// Don't forget to validate the alg is what you expect:
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			msg := fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
			return nil, msg
		}
		return []byte(sugar), nil
	})

	if err != nil {
		return nil, errors.Wrap(err, "Error parsing token")

	}

	if token == nil || !token.Valid {
		return nil, errors.New("invalid token")
	}

	user := model.User{
		ID:             uint(claims["id"].(float64)),
		Username:       claims["username"].(string),
		Role:           claims["role"].(string),
		BranchID:       uint(claims["branchId"].(float64)),
		OrganizationID: uint(claims["organizationId"].(float64)),
	}
	return &user, nil
}
