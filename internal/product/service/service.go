package service

import (
	"context"
	"gomicroservices/internal/product/model"
	"gomicroservices/internal/product/repo"

	"github.com/jackc/pgx/v4/pgxpool"
)

type Service struct {
	repo repo.Repo
}

func New(db *pgxpool.Pool) Service {
	return Service{
		repo: repo.New(db),
	}
}

func (s Service) CreateBrand(ctx context.Context, brand model.Brand) error {
	return s.repo.CreateBrand(ctx, brand)
}

func (s Service) GetBrand(ctx context.Context, id uint) (*model.Brand, error) {
	return s.repo.GetBrand(ctx, id)
}

func (s Service) GetBrands(ctx context.Context) ([]*model.Brand, error) {
	return s.repo.GetBrands(ctx)
}

func (s Service) CreateCategory(ctx context.Context, category model.Category) error {
	return s.repo.CreateCategory(ctx, category)
}

func (s Service) GetCategory(ctx context.Context, id uint) (*model.Category, error) {
	return s.repo.GetCategory(ctx, id)
}

func (s Service) GetCategories(ctx context.Context) ([]*model.Category, error) {
	return s.repo.GetCategories(ctx)
}

func (s Service) CreateProduct(ctx context.Context, product model.Product) error {
	return s.repo.CreateProduct(ctx, product)
}

func (s Service) GetProduct(ctx context.Context, id uint) (*model.Product, error) {
	return s.repo.GetProduct(ctx, id)
}

func (s Service) GetProducts(ctx context.Context) ([]*model.Product, error) {
	return s.repo.GetProducts(ctx)
}
