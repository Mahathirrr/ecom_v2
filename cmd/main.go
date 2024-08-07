package main

import (
	"ecom_apiv1/db"
	"ecom_apiv1/internal/handler"
	"ecom_apiv1/internal/server"
	"ecom_apiv1/internal/storer"
	"log"
	"os"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

const minSecretKeySize = 32

func main() {
	err := godotenv.Load("../.env")
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	secretKey := os.Getenv("SECRET_KEY")
	if len(secretKey) < minSecretKeySize {
		log.Fatalf("SECRET_KEY must be at least %d characters", minSecretKeySize)
	}
	gormDB, err := db.GetConnection()
	if err != nil {
		log.Fatalf("error opening database: %v", err)
	}
	log.Println("Succesfully connecting database")

	// str := storer.NewMySQLStorage(sqlx)
	str := storer.NewGORMStorage(gormDB)
	srv := server.NewServer(str)

	hdl := handler.NewHandler(srv, secretKey)
	handler.RegisterRoutes(hdl)
	handler.Start(":8000")
}
