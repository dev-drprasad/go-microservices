package service

import (
	"context"
	"gomicroservices/internal/order/model"
	"gomicroservices/internal/order/repo"
	"gomicroservices/internal/util"

	"github.com/jackc/pgx/v4/pgxpool"
)

// type IService interface {
// 	GetBranch(ctx context.Context, id uint) (*model.Branch, error)
// 	CreateOrganization(ctx context.Context, organization model.Organization) error
// 	CreateBranch(ctx context.Context, branch model.Branch) error
// 	GetOrganizations(ctx context.Context) ([]*model.Organization, error)
// }

type Service struct {
	repo repo.Repo
}

func New(db *pgxpool.Pool) Service {

	return Service{
		repo: repo.New(db),
	}
}

func (s Service) PlaceOrder(ctx context.Context, o model.Order) error {
	return s.repo.PlaceOrder(ctx, o)
}

func (s Service) GetOrder(ctx context.Context, id uint) (*model.Order, error) {
	log := util.GetLoggerFromContext(ctx)

	log.Infof("Reading order from repo. id=%d", id)
	return s.repo.GetOrder(ctx, id)
}

func (s Service) GetOrders(ctx context.Context) ([]*model.Order, error) {
	return s.repo.GetOrders(ctx)
}
