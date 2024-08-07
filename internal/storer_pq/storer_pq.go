package storerpq

import (
	"context"
	"ecom_apiv1/internal/storer"
	"fmt"

	"github.com/jmoiron/sqlx"
)

type PostgresStorage struct {
	*sqlx.DB
}

func NewPostgresStorage(db *sqlx.DB) *PostgresStorage {
	return &PostgresStorage{
		DB: db,
	}
}

// func (ps *PostgresStorage) CreateProduct(ctx context.Context, p *storer.Product) (*storer.Product, error) {
// 	res, err := ps.DB.NamedExecContext(ctx, "INSERT INTO products (name, image, category, description, rating, num_reviews, price, count_in_stock) VALUES (:name, :image, :category, :description, :rating, :num_reviews, :price, :count_in_stock)", p)
// 	if err != nil {
// 		return nil, fmt.Errorf("error inserting product: %w", err)
// 	}
//
// 	var id int
// 	err = ps.DB.GetContext(ctx, &id, "SELECT LASTVAL()")
// 	if err != nil {
// 		return nil, fmt.Errorf("error getting last insert ID: %w", err)
// 	}
// 	p.ID = id
// 	return p, nil
// }

func (ps *PostgresStorage) CreateProduct(ctx context.Context, p *storer.Product) (*storer.Product, error) {
	query := `
		INSERT INTO products (name, image, category, description, rating, num_reviews, price, count_in_stock) 
		VALUES (:name, :image, :category, :description, :rating, :num_reviews, :price, :count_in_stock) 
		RETURNING id`

	stmt, err := ps.DB.PrepareNamedContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("error preparing named statement for product: %w", err)
	}
	defer stmt.Close()

	err = stmt.GetContext(ctx, &p.ID, p)
	if err != nil {
		return nil, fmt.Errorf("error inserting product and getting ID: %w", err)
	}

	return p, nil
}

func (ps *PostgresStorage) GetProduct(ctx context.Context, id int) (*storer.Product, error) {
	var p storer.Product
	err := ps.DB.GetContext(ctx, &p, "SELECT * FROM products WHERE ID = $1", id)
	if err != nil {
		return nil, fmt.Errorf("error getting product: %w", err)
	}
	return &p, nil
}

func (ps *PostgresStorage) ListProducts(ctx context.Context) ([]storer.Product, error) {
	var products []storer.Product
	err := ps.DB.SelectContext(ctx, &products, "SELECT * FROM products")
	if err != nil {
		return nil, fmt.Errorf("error listing products: %w", err)
	}
	return products, nil
}

func (ps *PostgresStorage) UpdateProduct(ctx context.Context, p *storer.Product) (*storer.Product, error) {
	_, err := ps.DB.NamedExecContext(ctx, "UPDATE products SET name=:name, image=:image, category=:category, description=:description, rating=:rating, num_reviews=:num_reviews, price=:price, count_in_stock=:count_in_stock WHERE id=:id", p)
	if err != nil {
		return nil, fmt.Errorf("error updating product: %w", err)
	}
	return p, nil
}

func (ps *PostgresStorage) DeleteProduct(ctx context.Context, id int) error {
	_, err := ps.DB.ExecContext(ctx, "DELETE FROM products WHERE id=$1", id)
	if err != nil {
		return fmt.Errorf("error deleting product: %w", err)
	}
	return nil
}

func (ps *PostgresStorage) execTx(ctx context.Context, fn func(tx *sqlx.Tx) error) error {
	tx, err := ps.DB.BeginTxx(ctx, nil)
	if err != nil {
		return fmt.Errorf("error starting transaction: %w", err)
	}

	err = fn(tx)
	if err != nil {
		if rbErr := tx.Rollback(); rbErr != nil {
			return fmt.Errorf("error rolling back transaction: %w", err)
		}
		return fmt.Errorf("error in transaction: %w", err)
	}
	if err := tx.Commit(); err != nil {
		return fmt.Errorf("error committing transaction: %w", err)
	}

	return nil
}

func createOrder(ctx context.Context, tx *sqlx.Tx, o *storer.Order) (*storer.Order, error) {
	query := `
		INSERT INTO orders (payment_method, tax_price, shipping_price, total_price, user_id) 
		VALUES (:payment_method, :tax_price, :shipping_price, :total_price, :user_id) 
		RETURNING id`

	stmt, err := tx.PrepareNamedContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("error preparing named statement for order: %w", err)
	}
	defer stmt.Close()

	err = stmt.GetContext(ctx, &o.ID, o)
	if err != nil {
		return nil, fmt.Errorf("error inserting order and getting ID: %w", err)
	}

	return o, nil
}

func createOrderItem(ctx context.Context, tx *sqlx.Tx, oi *storer.OrderItem) error {
	_, err := tx.NamedExecContext(ctx, "INSERT INTO order_items (name, quantity, image, price, product_id, order_id) VALUES (:name, :quantity, :image, :price, :product_id, :order_id)", oi)
	if err != nil {
		return fmt.Errorf("error inserting order item: %w", err)
	}
	return nil
}

func (ps *PostgresStorage) CreateOrder(ctx context.Context, o *storer.Order) (*storer.Order, error) {
	err := ps.execTx(ctx, func(tx *sqlx.Tx) error {
		order, err := createOrder(ctx, tx, o)
		if err != nil {
			return fmt.Errorf("error creating order: %w", err)
		}

		for _, oi := range o.Items {
			oi.OrderID = order.ID
			err = createOrderItem(ctx, tx, &oi)
			if err != nil {
				return fmt.Errorf("error creating order item: %w", err)
			}
		}
		return nil
	})
	if err != nil {
		return nil, fmt.Errorf("error creating order: %w", err)
	}
	return o, nil
}

func (ps *PostgresStorage) GetOrder(ctx context.Context, userId int) (*storer.Order, error) {
	var o storer.Order
	err := ps.DB.GetContext(ctx, &o, "SELECT * FROM orders WHERE user_id=$1", userId)
	if err != nil {
		return nil, fmt.Errorf("error getting order: %w", err)
	}
	var oi []storer.OrderItem
	err = ps.DB.SelectContext(ctx, &oi, "SELECT * FROM order_items WHERE order_id=$1", o.ID)
	if err != nil {
		return nil, fmt.Errorf("error getting order items: %w", err)
	}
	o.Items = oi
	return &o, nil
}

func (ps *PostgresStorage) ListOrders(ctx context.Context) ([]storer.Order, error) {
	var orders []storer.Order
	err := ps.DB.SelectContext(ctx, &orders, "SELECT * FROM orders")
	if err != nil {
		return nil, fmt.Errorf("error listing orders: %w", err)
	}
	for i := range orders {
		var items []storer.OrderItem
		err := ps.DB.SelectContext(ctx, &items, "SELECT * FROM order_items WHERE order_id=$1", orders[i].ID)
		if err != nil {
			return nil, fmt.Errorf("error listing order items: %w", err)
		}
		orders[i].Items = items
	}
	return orders, nil
}

func (ps *PostgresStorage) DeleteOrder(ctx context.Context, id int) error {
	err := ps.execTx(ctx, func(tx *sqlx.Tx) error {
		_, err := tx.ExecContext(ctx, "DELETE FROM order_items WHERE order_id=$1", id)
		if err != nil {
			return fmt.Errorf("error deleting order items: %w", err)
		}

		_, err = tx.ExecContext(ctx, "DELETE FROM orders WHERE id=$1", id)
		if err != nil {
			return fmt.Errorf("error deleting order: %w", err)
		}
		return nil
	})
	if err != nil {
		return fmt.Errorf("error deleting order: %w", err)
	}
	return nil
}

func (ps *PostgresStorage) CreateUser(ctx context.Context, u *storer.User) (*storer.User, error) {
	query := `
		INSERT INTO users (name, email, password, is_admin) 
		VALUES (:name, :email, :password, :is_admin) 
		RETURNING id`

	stmt, err := ps.DB.PrepareNamedContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("error preparing named statement for user: %w", err)
	}
	defer stmt.Close()

	err = stmt.GetContext(ctx, &u.ID, u)
	if err != nil {
		return nil, fmt.Errorf("error inserting user and getting ID: %w", err)
	}

	return u, nil
}

func (ps *PostgresStorage) GetUser(ctx context.Context, email string) (*storer.User, error) {
	var u storer.User
	err := ps.DB.GetContext(ctx, &u, "SELECT * FROM users WHERE email=$1", email)
	if err != nil {
		return nil, fmt.Errorf("error getting user: %w", err)
	}

	return &u, nil
}

func (ps *PostgresStorage) ListUsers(ctx context.Context) ([]storer.User, error) {
	var users []storer.User
	err := ps.DB.SelectContext(ctx, &users, "SELECT * FROM users")
	if err != nil {
		return nil, fmt.Errorf("error listing users: %w", err)
	}

	return users, nil
}

func (ps *PostgresStorage) UpdateUser(ctx context.Context, u *storer.User) (*storer.User, error) {
	_, err := ps.DB.NamedExecContext(ctx, "UPDATE users SET name=:name, email=:email, password=:password, is_admin=:is_admin WHERE id=:id", u)
	if err != nil {
		return nil, fmt.Errorf("error updating user: %w", err)
	}

	return u, nil
}

func (ps *PostgresStorage) DeleteUser(ctx context.Context, id int) error {
	_, err := ps.DB.ExecContext(ctx, "DELETE FROM users WHERE id=$1", id)
	if err != nil {
		return fmt.Errorf("error deleting user: %w", err)
	}

	return nil
}

func (ps *PostgresStorage) CreateSession(ctx context.Context, s *storer.Session) (*storer.Session, error) {
	_, err := ps.DB.NamedExecContext(ctx, "INSERT INTO sessions (id, user_email, refresh_token, is_revoked, expires_at) VALUES (:id, :user_email, :refresh_token, :is_revoked, :expires_at)", s)
	if err != nil {
		return nil, fmt.Errorf("error inserting session: %w", err)
	}

	return s, nil
}

func (ps *PostgresStorage) GetSession(ctx context.Context, id string) (*storer.Session, error) {
	var s storer.Session
	err := ps.DB.GetContext(ctx, &s, "SELECT * FROM sessions WHERE id=$1", id)
	if err != nil {
		return nil, fmt.Errorf("error getting session: %w", err)
	}

	return &s, nil
}

func (ps *PostgresStorage) RevokeSession(ctx context.Context, id string) error {
	_, err := ps.DB.NamedExecContext(ctx, "UPDATE sessions SET is_revoked=1 WHERE id=:id", map[string]interface{}{"id": id})
	if err != nil {
		return fmt.Errorf("error revoking session: %w", err)
	}

	return nil
}

func (ps *PostgresStorage) DeleteSession(ctx context.Context, id string) error {
	_, err := ps.DB.ExecContext(ctx, "DELETE FROM sessions WHERE id=$1", id)
	if err != nil {
		return fmt.Errorf("error deleting session: %w", err)
	}

	return nil
}
