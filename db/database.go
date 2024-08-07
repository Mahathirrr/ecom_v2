package db

import (
	"fmt"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// func GetConnection() (*sqlx.DB, error) {
// 	db, err := sqlx.Open("postgres", "user=dev password=dev dbname=go_kasus_1 sslmode=disable")
// 	// db, err := sqlx.Open("mysql", "dev:dev@tcp(localhost:3306)/go_kasus_4?parseTime=true")
// 	if err != nil {
// 		return nil, fmt.Errorf("Error opening database: %w", err)
// 	}
// 	return db, nil
// }

func GetConnection() (*gorm.DB, error) {
	dsn := "dev:dev@tcp(localhost:3306)/go_kasus_5?charset=utf8mb4&parseTime=True&loc=Local"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		return nil, fmt.Errorf("Error opening database: %w", err)
	}
	return db, nil
}
