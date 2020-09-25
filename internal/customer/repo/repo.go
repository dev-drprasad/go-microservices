package repo

import (
	"context"
	"database/sql"
	"fmt"
	"gomicroservices/internal/customer/model"
	"gomicroservices/internal/util"
	"strings"

	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/pkg/errors"
)

var tableName = "customers"

type Repo interface {
	AddCustomer(ctx context.Context, user *model.Customer) error
	GetCustomer(ctx context.Context, id uint) (*model.Customer, error)
	UpdateCustomer(ctx context.Context, id uint, c model.Customer) error
	GetCustomers(ctx context.Context, organizationID uint) ([]*model.Customer, error)
	NewCustomersCount(ctx context.Context) ([]*util.CountOnDate, error)
}

type DBRepo struct {
	db *pgxpool.Pool
}

func New(db *pgxpool.Pool) DBRepo {
	return DBRepo{db: db}
}

func (repo DBRepo) GetCustomer(ctx context.Context, id uint) (*model.Customer, error) {

	stmt := `SELECT id, createdAt, updatedAt, name, address, zipcode, phoneNumber FROM %s WHERE id = $1`
	stmt = fmt.Sprintf(stmt, tableName)

	var c model.Customer
	err := repo.db.QueryRow(ctx, stmt, id).Scan(&c.ID, &c.CreatedAt, &c.UpdatedAt, &c.Name, &c.Address, &c.Zipcode, &c.PhoneNumber)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, util.ErrNoResourceFound
		}
		return nil, errors.Wrap(err, "Failed to execute the query")
	}

	return &c, nil
}

func (repo DBRepo) UpdateCustomer(ctx context.Context, id uint, c model.Customer) error {
	stmt := `
    UPDATE $1
    SET name = $2, address = $3, zipcode = $4, phoneNumber = $5
    WHERE id = $6
  `
	_, err := repo.db.Exec(ctx, stmt, tableName, c.Name, c.Address, c.Zipcode, c.PhoneNumber, id)
	if err != nil {
		if strings.Contains(err.Error(), util.ErrFKViolation.Error()) {
			return util.ErrFKViolation
		}
		return errors.Wrapf(err, "Failed to update customer. id=%v", id)
	}
	return nil
}

func (repo DBRepo) AddCustomer(ctx context.Context, c *model.Customer) error {
	stmt := `
    INSERT INTO %s
      (name, address, zipcode, phoneNumber)
    VALUES ($1, $2, $3, $4)
  `
	stmt = fmt.Sprintf(stmt, tableName)

	_, err := repo.db.Exec(ctx, stmt, c.Name, c.Address, c.Zipcode, c.PhoneNumber)
	if err != nil {
		if strings.Contains(err.Error(), util.ErrFKViolation.Error()) {
			return util.ErrFKViolation
		}
		return errors.Wrapf(err, "Failed to execute the query : %v", err)
	}

	return nil
}

func (repo DBRepo) GetCustomers(ctx context.Context) ([]*model.Customer, error) {

	customers := []*model.Customer{}

	stmt := `SELECT id, createdAt, updatedAt, name, address, zipcode, phoneNumber FROM %s`
	stmt = fmt.Sprintf(stmt, tableName)

	rows, err := repo.db.Query(ctx, stmt)
	if err != nil {
		return nil, errors.Wrap(err, "Failed to execute the query")
	}

	for rows.Next() {
		var c model.Customer
		if err = rows.Scan(&c.ID, &c.CreatedAt, &c.UpdatedAt, &c.Name, &c.Address, &c.Zipcode, &c.PhoneNumber); err != nil {
			break
		}
		customers = append(customers, &c)
	}
	if err != nil {
		return customers, errors.Wrap(err, "Failed to scan rows")
	}

	return customers, nil
}

func (repo DBRepo) NewCustomersCount(ctx context.Context) ([]*util.CountOnDate, error) {

	// https://stackoverflow.com/a/21008768/6748719
	stmt := `
    SELECT DATE(ds.series) as createdAt, COUNT(id)
    FROM %s
    RIGHT OUTER JOIN (
      SELECT
      GENERATE_SERIES((CURRENT_DATE - INTERVAL '14 days'), CURRENT_DATE, '1 day')
      AS series
    ) AS ds 
    ON DATE(createdAt) = ds.series
    GROUP BY DATE(ds.series)
    ORDER BY createdAt
  `
	// stmt := `
	// 	SELECT DATE(createdAt) as createdAt, COUNT(id)
	// 	FROM %s
	// 	WHERE createdAt > (CURRENT_DATE - INTERVAL '14 days')
	// 	GROUP BY DATE(createdAt)
	// 	ORDER BY createdAt
	// `
	stmt = fmt.Sprintf(stmt, tableName)

	rows, err := repo.db.Query(ctx, stmt)
	if err != nil {
		return nil, errors.Wrap(err, "Failed to execute the query")
	}

	counts := []*util.CountOnDate{}
	for rows.Next() {
		var c util.CountOnDate
		if err = rows.Scan(&c.Date, &c.Count); err != nil {
			break
		}
		counts = append(counts, &c)
	}
	if err != nil {
		return counts, errors.Wrap(err, "Failed to scan rows")
	}

	return counts, nil
}
