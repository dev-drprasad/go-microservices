package service

import (
	"context"
	"fmt"
	"gomicroservices/internal/customer/model"
	"gomicroservices/internal/customer/repo"
	"gomicroservices/internal/util"

	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/pkg/errors"
)

var ErrInvalidRequest = errors.New("invalid request")

type Service struct {
	repo repo.DBRepo
}

func New(db *pgxpool.Pool) Service {

	return Service{
		repo: repo.New(db),
	}
}

func (s *Service) AddCustomer(ctx context.Context, customer *model.Customer) error {

	if err := s.repo.AddCustomer(ctx, customer); err != nil {
		if err == util.ErrFKViolation {
			return ErrInvalidRequest
		}
		return fmt.Errorf("Failed to add customer : %v", err)
	}
	return nil
}

func (s *Service) GetCustomer(ctx context.Context, id uint) (*model.Customer, error) {
	log := util.GetLoggerFromContext(ctx)

	log.Infof("Reading customers from repo. id=%v", id)
	u, err := s.repo.GetCustomer(ctx, id)
	if err != nil {
		return nil, err
	}

	return u, nil
}

func (s Service) UpdateCustomer(ctx context.Context, id uint, p model.Customer) error {
	return s.repo.UpdateCustomer(ctx, id, p)
}

func (s *Service) GetCustomers(ctx context.Context) ([]*model.Customer, error) {
	return s.repo.GetCustomers(ctx)
}

func (s *Service) NewCustomersCount(ctx context.Context) ([]*util.CountOnDate, error) {
	return s.repo.NewCustomersCount(ctx)
}
