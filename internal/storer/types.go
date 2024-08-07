package storer

import (
	"time"
)

type Product struct {
	ID           uint `gorm:"primaryKey"`
	CreatedAt    time.Time
	UpdatedAt    time.Time
	Name         string  `gorm:"not null"`
	Image        string  `gorm:"not null"`
	Category     string  `gorm:"not null"`
	Description  string  `gorm:"type:text"`
	Rating       int     `gorm:"not null"`
	NumReviews   int     `gorm:"not null;default:0"`
	Price        float64 `gorm:"not null;type:decimal(10,2)"`
	CountInStock int     `gorm:"not null"`
}

type Order struct {
	ID            uint `gorm:"primaryKey"`
	CreatedAt     time.Time
	UpdatedAt     time.Time
	PaymentMethod string      `gorm:"not null"`
	TaxPrice      float64     `gorm:"not null;type:decimal(10,2)"`
	ShippingPrice float64     `gorm:"not null;type:decimal(10,2)"`
	TotalPrice    float64     `gorm:"not null;type:decimal(10,2)"`
	UserID        uint        `gorm:"not null"`
	User          User        `gorm:"foreignKey:UserID"`
	Items         []OrderItem `gorm:"foreignKey:OrderID"`
}

type OrderItem struct {
	ID        uint `gorm:"primaryKey"`
	CreatedAt time.Time
	UpdatedAt time.Time
	Name      string  `gorm:"not null"`
	Quantity  int     `gorm:"not null"`
	Image     string  `gorm:"not null"`
	Price     float64 `gorm:"not null;type:decimal(10,2)"`
	ProductID uint    `gorm:"not null"`
	OrderID   uint    `gorm:"not null"`
	Product   Product `gorm:"foreignKey:ProductID"`
	Order     Order   `gorm:"foreignKey:OrderID"`
}

type User struct {
	ID        uint `gorm:"primaryKey"`
	CreatedAt time.Time
	UpdatedAt time.Time
	Name      string `gorm:"not null"`
	Email     string `gorm:"not null;uniqueIndex"`
	Password  string `gorm:"not null"`
	IsAdmin   bool   `gorm:"not null;default:false"`
}

type Session struct {
	ID           string    `gorm:"primaryKey"`
	UserEmail    string    `gorm:"not null"`
	RefreshToken string    `gorm:"not null;type:varchar(512)"`
	IsRevoked    bool      `gorm:"not null;default:false"`
	CreatedAt    time.Time `gorm:"autoCreateTime"`
	ExpiresAt    time.Time `gorm:"not null"`
}
