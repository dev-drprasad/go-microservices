package service

import (
	"context"
	"database/sql"
	"gomicroservices/internal/organization/model"
	"gomicroservices/internal/organization/repo"
	dbrepo "gomicroservices/internal/organization/repo"
	"gomicroservices/internal/util"
)

type IService interface {
	GetBranch(ctx context.Context, id uint64) *model.Branch
}

type Service struct {
	repo repo.Repo
}

// var repoType string

// func init() {
// 	flag.StringVar(&repoType, "orgsvc", "db", "")
// }

func New(db *sql.DB) Service {
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
		repo: dbrepo.New(db),
	}
}

func (s Service) GetBranch(ctx context.Context, id uint64) *model.Branch {
	log := util.GetLoggerFromContext(ctx)

	log.Infof("Reading branch from repo. id=%d", id)
	return s.repo.GetBranch(ctx, id)
}
