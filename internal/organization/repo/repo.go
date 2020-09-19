package repo

import (
	"context"
	"database/sql"
	"gomicroservices/internal/organization/model"
)

type Repo interface {
	GetBranch(ctx context.Context, id uint64) *model.Branch
}

type DBRepo struct {
	db *sql.DB
}

func New(db *sql.DB) DBRepo {
	return DBRepo{db: db}
}

func (repo DBRepo) GetBranch(ctx context.Context, id uint64) *model.Branch {
	return &model.Branch{Name: "ABC #1", PhoneNumber: "0000000000"}
}
