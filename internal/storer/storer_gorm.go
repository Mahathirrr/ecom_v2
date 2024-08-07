package storer

import (
	"context"
	"errors"
	"fmt"

	"gorm.io/gorm"
)

var (
	ErrProductNotFound = errors.New("product not found")
	ErrUserNotFound    = errors.New("user not found")
	ErrOrderNotFound   = errors.New("order not found")
	ErrSessionNotFound = errors.New("session not found")
)

type GORMStorage struct {
	DB *gorm.DB
}

func NewGORMStorage(db *gorm.DB) *GORMStorage {
	return &GORMStorage{
		DB: db,
	}
}

func (gs *GORMStorage) CreateProduct(ctx context.Context, p *Product) (*Product, error) {
	result := gs.DB.WithContext(ctx).Create(p)
	if result.Error != nil {
		return nil, fmt.Errorf("error inserting product: %w", result.Error)
	}
	return p, nil
}

func (gs *GORMStorage) GetProduct(ctx context.Context, id uint) (*Product, error) {
	var p Product
	result := gs.DB.WithContext(ctx).First(&p, id)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, ErrProductNotFound
		}
		return nil, fmt.Errorf("error getting product: %w", result.Error)
	}
	return &p, nil
}

func (gs *GORMStorage) ListProducts(ctx context.Context) ([]Product, error) {
	var products []Product
	result := gs.DB.WithContext(ctx).Find(&products)
	if result.Error != nil {
		return nil, fmt.Errorf("error listing products: %w", result.Error)
	}
	return products, nil
}

func (gs *GORMStorage) UpdateProduct(ctx context.Context, p *Product) (*Product, error) {
	result := gs.DB.WithContext(ctx).Save(p)
	if result.Error != nil {
		return nil, fmt.Errorf("error updating product: %w", result.Error)
	}
	return p, nil
}

func (gs *GORMStorage) DeleteProduct(ctx context.Context, id uint) error {
	result := gs.DB.WithContext(ctx).Delete(&Product{}, id)
	if result.Error != nil {
		return fmt.Errorf("error deleting product: %w", result.Error)
	}
	if result.RowsAffected == 0 {
		return ErrProductNotFound
	}
	return nil
}

/*
func (gs *GORMStorage) CreateOrder(ctx context.Context, o *Order) (*Order, error) {
	err := gs.DB.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(o).Error; err != nil {
			return fmt.Errorf("error creating order: %w", err)
		}
		for i := range o.Items {
			o.Items[i].OrderID = o.ID
			if err := tx.Create(&o.Items[i]).Error; err != nil {
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
*/

// teknik bulk insert -> memasukkan data yang banyak sekaligus tanpa 1-1 ke db
func (gs *GORMStorage) CreateOrder(ctx context.Context, o *Order) (*Order, error) {
	err := gs.DB.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(o).Error; err != nil {
			return fmt.Errorf("error creating order: %w", err)
		}
		if len(o.Items) > 0 {
			for i := range o.Items {
				o.Items[i].OrderID = o.ID
			}
			if err := tx.Create(&o.Items).Error; err != nil {
				return fmt.Errorf("error creating order items: %w", err)
			}
		}
		return nil
	})
	if err != nil {
		return nil, fmt.Errorf("error creating order: %w", err)
	}
	return o, nil
}

func (gs *GORMStorage) GetOrder(ctx context.Context, userID uint) (*Order, error) {
	var o Order
	result := gs.DB.WithContext(ctx).Preload("Items").Where("user_id = ?", userID).First(&o)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, ErrOrderNotFound
		}
		return nil, fmt.Errorf("error getting order: %w", result.Error)
	}
	return &o, nil
}

func (gs *GORMStorage) ListOrders(ctx context.Context) ([]Order, error) {
	var orders []Order
	result := gs.DB.WithContext(ctx).Preload("Items").Find(&orders)
	if result.Error != nil {
		return nil, fmt.Errorf("error listing orders: %w", result.Error)
	}
	return orders, nil
}

func (gs *GORMStorage) DeleteOrder(ctx context.Context, id uint) error {
	err := gs.DB.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Where("order_id = ?", id).Delete(&OrderItem{}).Error; err != nil {
			return fmt.Errorf("error deleting order items: %w", err)
		}
		if err := tx.Delete(&Order{}, id).Error; err != nil {
			return fmt.Errorf("error deleting order: %w", err)
		}
		return nil
	})
	if err != nil {
		return fmt.Errorf("error deleting order: %w", err)
	}
	return nil
}

func (gs *GORMStorage) CreateUser(ctx context.Context, u *User) (*User, error) {
	result := gs.DB.WithContext(ctx).Create(u)
	if result.Error != nil {
		return nil, fmt.Errorf("error inserting user: %w", result.Error)
	}
	return u, nil
}

func (gs *GORMStorage) GetUser(ctx context.Context, email string) (*User, error) {
	var u User
	result := gs.DB.WithContext(ctx).Where("email = ?", email).First(&u)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, ErrUserNotFound
		}
		return nil, fmt.Errorf("error getting user: %w", result.Error)
	}
	return &u, nil
}

func (gs *GORMStorage) ListUsers(ctx context.Context) ([]User, error) {
	var users []User
	result := gs.DB.WithContext(ctx).Find(&users)
	if result.Error != nil {
		return nil, fmt.Errorf("error listing users: %w", result.Error)
	}
	return users, nil
}

func (gs *GORMStorage) UpdateUser(ctx context.Context, u *User) (*User, error) {
	result := gs.DB.WithContext(ctx).Save(u)
	if result.Error != nil {
		return nil, fmt.Errorf("error updating user: %w", result.Error)
	}
	return u, nil
}

func (gs *GORMStorage) DeleteUser(ctx context.Context, id uint) error {
	result := gs.DB.WithContext(ctx).Delete(&User{}, id)
	if result.Error != nil {
		return fmt.Errorf("error deleting user: %w", result.Error)
	}
	if result.RowsAffected == 0 {
		return ErrUserNotFound
	}
	return nil
}

func (gs *GORMStorage) CreateSession(ctx context.Context, s *Session) (*Session, error) {
	result := gs.DB.WithContext(ctx).Create(s)
	if result.Error != nil {
		return nil, fmt.Errorf("error inserting session: %w", result.Error)
	}
	return s, nil
}

func (gs *GORMStorage) GetSession(ctx context.Context, id string) (*Session, error) {
	var s Session
	result := gs.DB.WithContext(ctx).First(&s, "id = ?", id)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, ErrSessionNotFound
		}
		return nil, fmt.Errorf("error getting session: %w", result.Error)
	}
	return &s, nil
}

func (gs *GORMStorage) RevokeSession(ctx context.Context, id string) error {
	result := gs.DB.WithContext(ctx).Model(&Session{}).Where("id = ?", id).Update("is_revoked", true)
	if result.Error != nil {
		return fmt.Errorf("error revoking session: %w", result.Error)
	}
	if result.RowsAffected == 0 {
		return ErrSessionNotFound
	}
	return nil
}

func (gs *GORMStorage) DeleteSession(ctx context.Context, id string) error {
	result := gs.DB.WithContext(ctx).Delete(&Session{}, "id = ?", id)
	if result.Error != nil {
		return fmt.Errorf("error deleting session: %w", result.Error)
	}
	if result.RowsAffected == 0 {
		return ErrSessionNotFound
	}
	return nil
}
