package repo

import (
	"context"
	"fmt"
	customermodel "gomicroservices/internal/customer/model"
	"gomicroservices/internal/order/model"
	"gomicroservices/internal/util"
	"strings"

	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/pkg/errors"
)

const ordersTableName string = "orders"
const orderProductsTableName string = "order_products"

type Repo interface {
	PlaceOrder(ctx context.Context, branch model.Order) error
	GetOrder(ctx context.Context, id uint) (*model.Order, error)
	GetOrders(ctx context.Context) ([]*model.Order, error)
}

type DBRepo struct {
	db *pgxpool.Pool
}

func New(db *pgxpool.Pool) DBRepo {
	return DBRepo{db: db}
}

func (repo DBRepo) PlaceOrder(ctx context.Context, o model.Order) error {
	// stmt := `
	//   WITH newOrder AS (
	//     INSERT INTO %s (customerId) VALUES (?) returning id
	//   ) INSERT INTO %s (
	//     orderId, productId, unitPrice, quantity
	//   )
	//   VALUES %s
	// `
	// s := "UPDATE products SET stock = t.stock FROM (VALUES %s) AS t(id, stock) WHERE products.id = t.id"
	// s = util.MakeInsertSQLWithValue(s, "CAST(? AS int), (SELECT stock - ? FROM products WHERE id = CAST(? AS int))", len(o.Products), 1)
	// _, err := repo.db.Exec(ctx, s, 6, 1, 6)
	// fmt.Println("err", err)

	tx, err := repo.db.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return errors.Wrap(err, "transaction begin failed")
	}
	defer tx.Rollback(ctx)

	var orderID uint
	if err = tx.QueryRow(ctx, fmt.Sprintf("INSERT INTO %s (customerId) VALUES ($1) returning id", ordersTableName), o.CustomerID).Scan(&orderID); err != nil {
		return errors.Wrap(err, "order insertion failed")
	}

	stmt := "INSERT INTO order_products (orderId, productId, unitPrice, quantity) VALUES %s"
	stmt = util.MakeInsertSQLWithValue(stmt, "?, ?, (SELECT sellPrice FROM products WHERE id = ?) ,?", len(o.Products), 1)

	args := []interface{}{}
	for _, p := range o.Products {
		args = append(args, orderID, p.ProductID, p.ProductID, p.Quantity)
	}
	if _, err := tx.Exec(ctx, stmt, args...); err != nil {
		if strings.Contains(err.Error(), util.ErrFKViolation.Error()) {
			return util.ErrFKViolation
		}
		return errors.Wrap(err, "order insertion failed")
	}

	stmt = "UPDATE products SET stock = t.stock FROM (VALUES %s) AS t(id, stock) WHERE products.id = t.id"
	stmt = util.MakeInsertSQLWithValue(stmt, "CAST(? AS int), (SELECT stock - ? FROM products WHERE id = CAST(? AS int))", len(o.Products), 1)

	args = []interface{}{}
	for _, p := range o.Products {
		args = append(args, p.ProductID, p.Quantity, p.ProductID)
	}
	fmt.Println("args", args)
	if _, err := tx.Exec(ctx, stmt, args...); err != nil {
		return errors.Wrap(err, "product stock decrement failed")
	}

	if err = tx.Commit(ctx); err != nil {
		return errors.Wrap(err, "tx commit failed")
	}

	return nil
}

func (repo DBRepo) GetOrder(ctx context.Context, id uint) (*model.Order, error) {
	stmt := `
    SELECT o.id, o.customerId, op.productId, op.unitPrice, op.quantity
    FROM orders AS o
    LEFT JOIN order_products AS op ON o.id = op.orderId
    WHERE id = $1
  `

	rows, err := repo.db.Query(ctx, stmt, id)
	if err != nil {
		return nil, errors.Wrap(err, "Failed to execute the query")
	}
	var order model.Order
	for rows.Next() {
		var op model.OrderProduct
		if err = rows.Scan(&order.ID, &order.CustomerID, &op.ProductID, &op.UnitPrice, &op.Quantity); err != nil {
			break
		}
		order.Products = append(order.Products, &op)
	}

	if err != nil {
		return nil, errors.Wrap(err, "Failed to execute the query")
	}

	return &order, nil
}

func (repo DBRepo) GetOrders(ctx context.Context) ([]*model.Order, error) {
	stmt := `
    SELECT o.id, o.status, c.id, c.name, JSON_AGG(JSON_BUILD_OBJECT('productId', op.productId, 'quantity', op.quantity, 'unitPrice', op.unitPrice))
    FROM orders AS o
    LEFT JOIN order_products AS op ON o.id = op.orderId
    JOIN customers AS c ON o.customerId = c.id
    GROUP BY o.id, c.id
  `

	orders := []*model.Order{}
	rows, err := repo.db.Query(ctx, stmt)
	if err != nil {
		return nil, errors.Wrap(err, "Failed to execute the query")
	}

	for rows.Next() {
		var o model.Order
		var ops []*model.OrderProduct
		var c customermodel.Customer
		if err = rows.Scan(&o.ID, &o.Status, &c.ID, &c.Name, &ops); err != nil {
			break
		}
		o.Products = ops
		o.Customer = &c
		orders = append(orders, &o)
	}
	if err != nil {
		return orders, errors.Wrap(err, "Failed to scan rows")
	}

	return orders, nil
}
