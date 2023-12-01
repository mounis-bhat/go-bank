package main

import (
	"fmt"
	"log"

	"github.com/joho/godotenv"
)

func init() {
	if err := godotenv.Load(); err != nil {
		fmt.Println("Error loading .env file")
	}
}

func main() {
	store, err := NewPostgresStorage()
	if err != nil {
		log.Fatal(err)
	}
	if err := InitStorage(store); err != nil {
		log.Fatal(err)
	}

	server := NewAPIServer(":8080", store)
	server.Run()
}
