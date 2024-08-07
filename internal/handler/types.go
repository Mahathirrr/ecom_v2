package handler

import "time"

type ProductReq struct {
	Name         string  `json:"name" validate:"required,min=3,max=255"`
	Image        string  `json:"image" validate:"required,url"`
	Category     string  `json:"category" validate:"required,min=2,max=255"`
	Description  string  `json:"description" validate:"max=1000"`
	Rating       int     `json:"rating" validate:"min=0,max=5"`
	NumReviews   int     `json:"num_reviews" validate:"min=0"`
	Price        float64 `json:"price" validate:"required,gt=0"`
	CountInStock int     `json:"count_in_stock" validate:"min=0"`
}

type ProductRes struct {
	ID           uint      `json:"id"`
	Name         string    `json:"name"`
	Image        string    `json:"image"`
	Category     string    `json:"category"`
	Description  string    `json:"description"`
	Rating       int       `json:"rating"`
	NumReviews   int       `json:"num_reviews"`
	Price        float64   `json:"price"`
	CountInStock int       `json:"count_in_stock"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at,omitempty"`
}

type OrderReq struct {
	Items         []OrderItem `json:"items" validate:"required,min=1,dive"`
	PaymentMethod string      `json:"payment_method" validate:"required,oneof=PayPal Stripe"`
	TaxPrice      float64     `json:"tax_price" validate:"min=0"`
	ShippingPrice float64     `json:"shipping_price" validate:"min=0"`
	TotalPrice    float64     `json:"total_price" validate:"required,gt=0"`
}

type OrderItem struct {
	Name      string  `json:"name"`
	Quantity  int     `json:"quantity"`
	Image     string  `json:"image"`
	Price     float64 `json:"price"`
	ProductID uint    `json:"product_id"`
}

type OrderRes struct {
	ID            uint        `json:"id"`
	Items         []OrderItem `json:"items"`
	PaymentMethod string      `json:"payment_method"`
	TaxPrice      float64     `json:"tax_price"`
	ShippingPrice float64     `json:"shipping_price"`
	TotalPrice    float64     `json:"total_price"`
	CreatedAt     time.Time   `json:"created_at"`
	UpdatedAt     time.Time   `json:"updated_at,omitempty"`
}

type UserReq struct {
	Name     string `json:"name" validate:"required,min=3,max=255"`
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=8"`
	IsAdmin  bool   `json:"is_admin"`
}

type UserRes struct {
	Name    string `json:"name"`
	Email   string `json:"email"`
	IsAdmin bool   `json:"is_admin"`
}

type ListUserRes struct {
	Users []UserRes `json:"users"`
}

type LoginUserReq struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=8"`
}

type LoginUserRes struct {
	SessionID             string    `json:"session_id"`
	AccessToken           string    `json:"access_token"`
	RefreshToken          string    `json:"refresh_token"`
	AccessTokenExpiresAt  time.Time `json:"access_token_expires_at"`
	RefreshTokenExpiresAt time.Time `json:"refresh_token_expires_at"`
	User                  UserRes   `json:"user"`
}

type RenewAccessTokenReq struct {
	RefreshToken string `json:"refresh_token" validate:"required"`
}

type RenewAccessTokenRes struct {
	AccessToken          string    `json:"access_token"`
	AccessTokenExpiresAt time.Time `json:"access_token_expires_at"`
}
