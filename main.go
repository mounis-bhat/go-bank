package main

import (
	"fmt"
)

func main() {
	server := NewAPIServer(":8080")
	server.Start()
	fmt.Println("Yeet! 🚀")
}
