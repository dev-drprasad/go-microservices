package repo

import (
	"context"
	"gomicroservices/internal/organization/model"
	"gomicroservices/internal/util"
	"strings"

	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/pkg/errors"
)

type Repo interface {
	GetBranches(ctx context.Context) ([]*model.Branch, error)
	GetBranch(ctx context.Context, id uint) (*model.Branch, error)
	CreateOrganization(ctx context.Context, organization model.Organization) error
	CreateBranch(ctx context.Context, branch model.Branch) error
	GetBranchesByOrganization(ctx context.Context, organizationID uint) ([]*model.Branch, error)
	GetOrganizations(ctx context.Context) ([]*model.Organization, error)
}

type DBRepo struct {
	db *pgxpool.Pool
}

func New(db *pgxpool.Pool) DBRepo {
	return DBRepo{db: db}
}

func (repo DBRepo) CreateOrganization(ctx context.Context, organization model.Organization) error {
	stmt := `INSERT INTO organizations (name) VALUES ($1)`

	_, err := repo.db.Exec(ctx, stmt, organization.Name)
	if err != nil {
		if strings.Contains(err.Error(), util.ErrFKViolation.Error()) {
			return util.ErrFKViolation
		}
		return errors.Wrapf(err, "Failed to execute the query name=%v", organization.Name)
	}

	return nil
}

func (repo DBRepo) CreateBranch(ctx context.Context, branch model.Branch) error {
	stmt := `INSERT INTO branches (name, address, phoneNumber, organizationId) VALUES ($1, $2, $3, $4)`

	_, err := repo.db.Exec(ctx, stmt, branch.Name, branch.Address, branch.PhoneNumber, branch.OrganizationID)
	if err != nil {
		if strings.Contains(err.Error(), util.ErrFKViolation.Error()) {
			return util.ErrFKViolation
		}
		return errors.Wrapf(err, "Failed to execute the query name=%v organizationId=%v", branch.Name, branch.OrganizationID)
	}

	return nil
}

func (repo DBRepo) GetBranch(ctx context.Context, id uint) (*model.Branch, error) {
	stmt := `SELECT id, name, organizationId FROM branches WHERE id = $1`

	var branch model.Branch
	err := repo.db.QueryRow(ctx, stmt, id).Scan(&branch.ID, &branch.Name, &branch.OrganizationID)
	if err != nil {
		return nil, errors.Wrap(err, "Failed to execute the query")
	}

	return &branch, nil
}

func (repo DBRepo) GetBranches(ctx context.Context) ([]*model.Branch, error) {
	stmt := `SELECT branches.id, branches.name FROM branches`

	var branches []*model.Branch
	rows, err := repo.db.Query(ctx, stmt)
	if err != nil {
		return nil, errors.Wrap(err, "Failed to execute the query")
	}

	for rows.Next() {
		var branch model.Branch
		if err = rows.Scan(&branch.ID, &branch.Name); err != nil {
			break
		}
		branches = append(branches, &branch)
	}
	if err != nil {
		return branches, errors.Wrap(err, "Failed to scan rows")
	}

	return branches, nil
}

func (repo DBRepo) GetBranchesByOrganization(ctx context.Context, organizationID uint) ([]*model.Branch, error) {
	stmt := `SELECT branches.id, branches.name FROM branches LEFT JOIN organizations ON organizations.id = branches.organizationId WHERE organizations.id = $1`

	var branches []*model.Branch
	rows, err := repo.db.Query(ctx, stmt, organizationID)
	if err != nil {
		return nil, errors.Wrap(err, "Failed to execute the query")
	}

	for rows.Next() {
		var branch model.Branch
		if err = rows.Scan(&branch.ID, &branch.Name); err != nil {
			break
		}
		branches = append(branches, &branch)
	}
	if err != nil {
		return branches, errors.Wrap(err, "Failed to scan rows")
	}

	return branches, nil
}

func (repo DBRepo) GetOrganizations(ctx context.Context) ([]*model.Organization, error) {
	stmt := `SELECT id, name FROM organizations`

	var organizations []*model.Organization
	rows, err := repo.db.Query(ctx, stmt)
	if err != nil {
		return nil, errors.Wrap(err, "Failed to execute the query")
	}

	for rows.Next() {
		var organization model.Organization
		if err = rows.Scan(&organization.ID, &organization.Name); err != nil {
			break
		}
		organizations = append(organizations, &organization)
	}
	if err != nil {
		return organizations, errors.Wrap(err, "Failed to scan rows")
	}

	return organizations, nil
}
