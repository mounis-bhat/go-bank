package main

import (
	"fmt"
	"log"

	"github.com/joho/godotenv"
	"github.com/mounis-bhat/go-bank/pkg/db"
	"github.com/mounis-bhat/go-bank/pkg/server"
)

func init() {
	if err := godotenv.Load(); err != nil {
		fmt.Println("Error loading .env file")
	}
}

func main() {
	store, err := db.NewPostgresStorage()
	if err != nil {
		log.Fatal(err)
	}
	if err := db.InitStorage(store); err != nil {
		log.Fatal(err)
	}

	server := server.NewAPIServer(":8080", store)
	server.Run()
}
