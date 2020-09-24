package service

import (
	"context"
	"gomicroservices/internal/organization/model"
	"gomicroservices/internal/organization/repo"
	usermodel "gomicroservices/internal/user/model"
	"gomicroservices/internal/util"

	"github.com/jackc/pgx/v4/pgxpool"
)

type IService interface {
	GetBranch(ctx context.Context, id uint) (*model.Branch, error)
	CreateOrganization(ctx context.Context, organization model.Organization) error
	CreateBranch(ctx context.Context, branch model.Branch) error
	GetOrganizations(ctx context.Context) ([]*model.Organization, error)
}

type Service struct {
	repo repo.Repo
}

// var repoType string

// func init() {
// 	flag.StringVar(&repoType, "orgsvc", "db", "")
// }

func New(db *pgxpool.Pool) Service {
	// var repo repo.Repo
	// switch repoType {
	// case "http":
	// 	repo = httprepo.HTTPRepo{}
	// case "grpc":
	// 	repo = grpc.GRPCRepo{}
	// default:
	// 	repo = dbrepo.DBRepo{}
	// }

	return Service{
		repo: repo.New(db),
	}
}

func (s Service) CreateOrganization(ctx context.Context, organization model.Organization) error {
	return s.repo.CreateOrganization(ctx, organization)
}

func (s Service) CreateBranch(ctx context.Context, branch model.Branch) error {
	return s.repo.CreateBranch(ctx, branch)
}

func (s Service) GetBranch(ctx context.Context, id uint) (*model.Branch, error) {
	log := util.GetLoggerFromContext(ctx)

	log.Infof("Reading branch from repo. id=%d", id)
	return s.repo.GetBranch(ctx, id)
}

func (s Service) GetBranchesByOrganization(ctx context.Context, organizationId uint) ([]*model.Branch, error) {
	log := util.GetLoggerFromContext(ctx)

	log.Infof("Reading branches of organization. organizationId=%d", organizationId)
	return s.repo.GetBranchesByOrganization(ctx, organizationId)
}

func (s Service) GetBranches(ctx context.Context) (branches []*model.Branch, err error) {
	log := util.GetLoggerFromContext(ctx)

	authUser := ctx.Value("user").(*usermodel.User)

	if authUser.Role == "superadmin" {
		branches, err = s.repo.GetBranches(ctx)
	} else {
		log.Infof("Reading branches of organization. organizationId=%d", authUser.OrganizationID)
		branches, err = s.repo.GetBranchesByOrganization(ctx, authUser.OrganizationID)
	}

	return
}

func (s Service) GetOrganizations(ctx context.Context) ([]*model.Organization, error) {

	return s.repo.GetOrganizations(ctx)
}
