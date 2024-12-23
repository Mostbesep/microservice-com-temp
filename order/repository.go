package order

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/lib/pq"
)

type Repository interface {
	Close() error
	PutOrder(ctx context.Context, order Order) error
	GetOrder(ctx context.Context, id string) (Order, error)
	GetAccountOrders(ctx context.Context, accountId string) (*[]Order, error)
}

type postgresqlRepository struct {
	db *sql.DB
}

func (r *postgresqlRepository) Close() error {
	err := r.db.Close()
	if err != nil {
		return err
	}
	return nil
}

func (r *postgresqlRepository) PutOrder(ctx context.Context, order Order) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer func() {
		if err != nil {
			tx.Rollback()
			return
		}
		err = tx.Commit()
	}()
	_, err = tx.ExecContext(
		ctx,
		"INSERT INTO orders (id, created_at, account_id, total_price) VALUES ($1, $2, $3, $4)",
		order.Id,
		order.CreatedAt,
		order.AccountId,
		order.TotalPrice,
	)
	if err != nil {
		return err
	}
	stmt, err := tx.PrepareContext(
		ctx, pq.CopyIn(
			"order_products",
			"order_id", "product_id", "quantity", "price"))
	if err != nil {
		return err
	}
	for _, p := range order.Products {
		_, err = stmt.ExecContext(ctx, order.AccountId, p.Id, p.Quantity)
		if err != nil {
			return err
		}
	}
	_, err = stmt.ExecContext(ctx)
	if err != nil {
		return err
	}
	err = stmt.Close()
	return err
}

func (r *postgresqlRepository) GetOrder(ctx context.Context, id string) (Order, error) {
	// Query the database
	rows, err := r.db.QueryContext(
		ctx, `
		SELECT
			o.id,
			o.created_at,
			o.account_id,
			o.total_price::float8,
			op.product_id,
			op.quantity
		FROM orders o
		JOIN order_products op ON (o.id = op.order_id)
		WHERE o.id = $1`,
		id)
	if err != nil {
		return Order{}, fmt.Errorf("failed to query order: %w", err)
	}
	defer rows.Close()

	var order Order
	var products []OrderedProduct
	orderFetched := false

	// Process rows
	for rows.Next() {
		var product OrderedProduct
		if err := rows.Scan(
			&order.Id,
			&order.CreatedAt,
			&order.Id,
			&order.TotalPrice,
			&product.Id,
			product.Quantity,
		); err != nil {
			return Order{}, fmt.Errorf("failed to scan row: %w", err)
		}
		products = append(products, product)
		orderFetched = true
	}

	// Check for iteration errors
	if err := rows.Err(); err != nil {
		return Order{}, fmt.Errorf("row iteration error: %w", err)
	}

	// Handle empty result set
	if !orderFetched {
		return Order{}, fmt.Errorf("order with ID %s not found", id)
	}

	order.Products = products
	return order, nil

}

func (r *postgresqlRepository) GetAccountOrders(ctx context.Context, accountId string) (*[]Order, error) {
	rows, err := r.db.QueryContext(
		ctx, `
	SELECT
	o.id, 
	o.created_at, 
	o.account_id, 
	o.total_price::money::numeric::float8, 
	op.product_id, 
	op.quantity 
	FROM orders o JOIN order_products op ON (o.id = op.order_id) 
	WHERE o.account_id = $1
	ORDER BY o.id`,
		accountId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	orders := []Order{}
	order := &Order{}
	lastOrder := &Order{}
	orderedProduct := &OrderedProduct{}
	products := []OrderedProduct{}
	// Scan rows into Order structs
	for rows.Next() {
		if err = rows.Scan(
			&order.Id,
			&order.CreatedAt,
			&order.AccountId,
			&order.TotalPrice,
			&orderedProduct.Id,
			&orderedProduct.Quantity,
		); err != nil {
			return nil, err
		}
		// Scan order
		if lastOrder.Id != "" && lastOrder.Id != order.Id {
			newOrder := Order{
				Id:         lastOrder.Id,
				AccountId:  lastOrder.AccountId,
				CreatedAt:  lastOrder.CreatedAt,
				TotalPrice: lastOrder.TotalPrice,
				Products:   products,
			}
			orders = append(orders, newOrder)
			products = []OrderedProduct{}
		}
		// Scan products
		products = append(products, OrderedProduct{
			Id:       orderedProduct.Id,
			Quantity: orderedProduct.Quantity,
		})

		*lastOrder = *order
	}

	// Add last order (or first :D)
	if lastOrder != nil {
		newOrder := Order{
			Id:         lastOrder.Id,
			AccountId:  lastOrder.AccountId,
			CreatedAt:  lastOrder.CreatedAt,
			TotalPrice: lastOrder.TotalPrice,
			Products:   products,
		}
		orders = append(orders, newOrder)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return &orders, nil

}

func NewPostgresqlRepository(url string) (Repository, error) {
	db, err := sql.Open("postgres", url)
	if err != nil {
		return nil, err
	}
	err = db.Ping()
	if err != nil {
		return nil, err
	}
	return &postgresqlRepository{db: db}, nil
}
