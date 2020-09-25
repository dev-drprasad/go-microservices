package repo

import (
	"context"
	"fmt"
	"gomicroservices/internal/product/model"
	"gomicroservices/internal/util"
	"strings"

	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/pkg/errors"
)

type Repo interface {
	CreateBrand(ctx context.Context, branch model.Brand) error
	GetBrand(ctx context.Context, id uint) (*model.Brand, error)
	UpdateProduct(ctx context.Context, id uint, p model.Product) error
	GetBrands(ctx context.Context) ([]*model.Brand, error)
	GetCategories(ctx context.Context) ([]*model.Category, error)
	GetCategory(ctx context.Context, id uint) (*model.Category, error)
	CreateCategory(ctx context.Context, branch model.Category) error
	CreateProduct(ctx context.Context, p model.Product) error
	GetProduct(ctx context.Context, id uint) (*model.Product, error)
	GetProducts(ctx context.Context) ([]*model.Product, error)
	AddProductImageURLs(ctx context.Context, imageURLs []string) error
}

type DBRepo struct {
	db *pgxpool.Pool
}

func New(db *pgxpool.Pool) DBRepo {
	return DBRepo{db: db}
}

func (repo DBRepo) CreateBrand(ctx context.Context, brand model.Brand) error {
	stmt := `INSERT INTO brands (name) VALUES ($1)`

	_, err := repo.db.Exec(ctx, stmt, brand.Name)
	if err != nil {
		if strings.Contains(err.Error(), util.ErrFKViolation.Error()) {
			return util.ErrFKViolation
		}
		return errors.Wrapf(err, "Failed to execute the query name=%v", brand.Name)
	}

	return nil
}

func (repo DBRepo) GetBrand(ctx context.Context, id uint) (*model.Brand, error) {
	stmt := `SELECT id, name FROM brands WHERE id = $1`

	var brand model.Brand
	err := repo.db.QueryRow(ctx, stmt, id).Scan(&brand.ID, &brand.Name)
	if err != nil {
		return nil, errors.Wrap(err, "Failed to execute the query")
	}

	return &brand, nil
}

func (repo DBRepo) GetBrands(ctx context.Context) ([]*model.Brand, error) {
	stmt := `SELECT id, name FROM brands`

	brands := []*model.Brand{}
	rows, err := repo.db.Query(ctx, stmt)
	if err != nil {
		return nil, errors.Wrap(err, "Failed to execute the query")
	}

	for rows.Next() {
		var brand model.Brand
		if err = rows.Scan(&brand.ID, &brand.Name); err != nil {
			break
		}
		brands = append(brands, &brand)
	}
	if err != nil {
		return brands, errors.Wrap(err, "Failed to scan rows")
	}

	return brands, nil
}

func (repo DBRepo) CreateCategory(ctx context.Context, category model.Category) error {
	stmt := `INSERT INTO categories (name) VALUES ($1)`

	_, err := repo.db.Exec(ctx, stmt, category.Name)
	if err != nil {
		if strings.Contains(err.Error(), util.ErrFKViolation.Error()) {
			return util.ErrFKViolation
		}
		return errors.Wrapf(err, "Failed to execute the query name=%v", category.Name)
	}

	return nil
}

func (repo DBRepo) GetCategory(ctx context.Context, id uint) (*model.Category, error) {
	stmt := `SELECT id, name FROM categories WHERE id = $1`

	var category model.Category
	err := repo.db.QueryRow(ctx, stmt, id).Scan(&category.ID, &category.Name)
	if err != nil {
		return nil, errors.Wrap(err, "Failed to execute the query")
	}

	return &category, nil
}

func (repo DBRepo) GetCategories(ctx context.Context) ([]*model.Category, error) {
	stmt := `SELECT id, name FROM categories`

	categories := []*model.Category{}
	rows, err := repo.db.Query(ctx, stmt)
	if err != nil {
		return nil, errors.Wrap(err, "Failed to execute the query")
	}

	for rows.Next() {
		var category model.Category
		if err = rows.Scan(&category.ID, &category.Name); err != nil {
			break
		}
		categories = append(categories, &category)
	}
	if err != nil {
		return categories, errors.Wrap(err, "Failed to scan rows")
	}

	return categories, nil
}

func (repo DBRepo) CreateProduct(ctx context.Context, p model.Product) error {
	stmt := `
		INSERT INTO products
			(name, cost, sellPrice, brandId, categoryId, imageUrls, stock)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
	`

	_, err := repo.db.Exec(ctx, stmt, p.Name, p.Cost, p.SellPrice, p.BrandID, p.CategoryID, p.ImageURLs, p.Stock)
	if err != nil {
		if strings.Contains(err.Error(), util.ErrFKViolation.Error()) {
			return util.ErrFKViolation
		}
		return errors.Wrapf(err, "Failed to execute the query name=%v", p.Name)
	}

	return nil
}

func (repo DBRepo) GetProduct(ctx context.Context, id uint) (*model.Product, error) {
	stmt := `
		SELECT
			id, name, cost, sellPrice, brandId, categoryId, imageUrls, stock
		FROM products
		WHERE id = $1
	`

	var p model.Product
	err := repo.db.QueryRow(ctx, stmt, id).Scan(&p.ID, &p.Name, &p.Cost, &p.SellPrice, &p.BrandID, &p.CategoryID, &p.ImageURLs, &p.Stock)
	if err != nil {
		return nil, errors.Wrap(err, "Failed to execute the query")
	}

	return &p, nil
}

func (repo DBRepo) UpdateProduct(ctx context.Context, id uint, p model.Product) error {
	stmt := `
		UPDATE products
		SET name = $1, cost = $2, sellPrice = $3, brandId = $4, categoryId = $5, imageUrls = $6, stock = $7
		WHERE id = $8
	`
	_, err := repo.db.Exec(ctx, stmt, p.Name, p.Cost, p.SellPrice, p.BrandID, p.CategoryID, p.ImageURLs, p.Stock, id)
	if err != nil {
		if strings.Contains(err.Error(), util.ErrFKViolation.Error()) {
			return util.ErrFKViolation
		}
		return errors.Wrapf(err, "Failed to update product. id=%v", id)
	}
	return nil
}

func (repo DBRepo) GetProducts(ctx context.Context) ([]*model.Product, error) {
	stmt := `
		SELECT
			id, name, cost, sellPrice, brandId, categoryId, imageUrls, stock
		FROM products
	`

	products := []*model.Product{}
	rows, err := repo.db.Query(ctx, stmt)
	if err != nil {
		return nil, errors.Wrap(err, "Failed to execute the query")
	}

	for rows.Next() {
		var p model.Product

		if err = rows.Scan(&p.ID, &p.Name, &p.Cost, &p.SellPrice, &p.BrandID, &p.CategoryID, &p.ImageURLs, &p.Stock); err != nil {
			break
		}

		products = append(products, &p)
	}
	if err != nil {
		return products, errors.Wrap(err, "Failed to scan rows")
	}

	return products, nil
}

func makeBulkInsertQuery(stmt string, perRow int, rows int) string {
	placeholders := make([]string, 0, rows)
	for i := 0; i < rows; i++ {
		placeholder := make([]string, 0, perRow)
		for j := i; j < i+perRow; j++ {
			placeholder = append(placeholder, fmt.Sprintf("$%d", j+i+1))
		}
		placeholderStr := strings.Join(placeholder, ", ")
		placeholders = append(placeholders, "("+placeholderStr+")")

	}
	return fmt.Sprintf(stmt, strings.Join(placeholders, ","))
}

func strtointerface(in []string) []interface{} {
	out := make([]interface{}, len(in))
	for i, v := range in {
		out[i] = v
	}
	return out
}

func (repo DBRepo) AddProductImageURLs(ctx context.Context, imageURLs []string) error {
	stmt := `
		INSERT INTO product_images
			(imageUrl)
		VALUES %s
	`
	stmt = makeBulkInsertQuery(stmt, 1, len(imageURLs))

	_, err := repo.db.Exec(ctx, stmt, strtointerface(imageURLs)...)
	if err != nil {
		if strings.Contains(err.Error(), util.ErrFKViolation.Error()) {
			return util.ErrFKViolation
		}
		return errors.Wrapf(err, "Failed to execute the query")
	}

	return nil
}
