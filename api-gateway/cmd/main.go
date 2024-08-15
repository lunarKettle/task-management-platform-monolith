package main

import (
	"api-gateway/internal/handler"
	"fmt"
)

func main() {
	server := handler.NewServer(":8080")
	fmt.Printf("Starting server at %s\n", server.Address)
	if err := server.Start(); err == nil {
		fmt.Println("Ошибка")
	}
}
