package repo

import (
	"context"
	"database/sql"
	"gomicroservices/internal/user/model"
	"gomicroservices/internal/util"

	"github.com/pkg/errors"
)

type Repo interface {
	GetUser(ctx context.Context, id uint) (*model.User, error)
	CreateUser(ctx context.Context, user *model.User) error
}

type DBRepo struct {
	db *sql.DB
}

func New(db *sql.DB) DBRepo {
	return DBRepo{db: db}
}

func (repo DBRepo) GetUser(ctx context.Context, id uint) (*model.User, error) {

	stmt := `SELECT id, name, username FROM users`
	var user model.User
	err := repo.db.QueryRow(stmt).Scan(&user.ID, &user.Name, &user.Username)
	if err != nil {
		return nil, errors.Wrap(err, "Failed to execute the query")
	}

	return &user, nil
}

func (repo DBRepo) CreateUser(ctx context.Context, user *model.User) error {
	log := util.GetLoggerFromContext(ctx)
	stmt := `INSERT INTO users (name, username, password) VALUES ($1, $2, crypt($3, gen_salt('bf'))) RETURNING id`
	var userID uint64

	err := repo.db.QueryRow(stmt, user.Name, user.Username, user.Password).Scan(&userID)
	if err != nil {
		return errors.Wrapf(err, "Failed to execute the query name=%v username=%v", user.Name, user.Username)
	}

	log.Infof("New user created with id=%v", userID)

	return nil
}
