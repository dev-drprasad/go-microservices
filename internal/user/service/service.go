package service

import (
	"context"
	"database/sql"
	"flag"
	orgsvc "gomicroservices/internal/organization/service"
	orgsvcgrpc "gomicroservices/internal/user/adapters/organization/grpc"
	"gomicroservices/internal/user/model"
	"gomicroservices/internal/user/repo"
	"gomicroservices/internal/util"
)

type Service struct {
	repo   repo.DBRepo
	orgsvc orgsvc.IService
}

var repoType string

func init() {
	flag.StringVar(&repoType, "orgsvc", "db", "")
	flag.Parse()
}

func New(db *sql.DB) Service {
	var o orgsvc.IService
	switch repoType {
	case "grpc":
		o = orgsvcgrpc.New()
	default:
		o = orgsvc.New(db)
	}
	return Service{
		repo:   repo.New(db),
		orgsvc: o,
	}
}

func (s *Service) GetUser(ctx context.Context, id uint) (*model.User, error) {
	log := util.GetLoggerFromContext(ctx)

	log.Infof("Reading user from repo. id=%s", id)
	u, err := s.repo.GetUser(ctx, id)
	if err != nil {
		return nil, err
	}

	log.Infof("Reading branch from org service. id=%s", u.BranchID)
	u.WorksAt = s.orgsvc.GetBranch(ctx, u.BranchID)
	return u, nil
}

func (s *Service) CreateUser(ctx context.Context, user *model.User) error {
	log := util.GetLoggerFromContext(ctx)

	log.Infof("Creating user with name=%v, username=%v", user.Name, user.Username)
	return s.repo.CreateUser(ctx, user)
}
