package main

import (
	"api_gateway/internal/grpc_client"
	"api_gateway/internal/handler"
	"fmt"
)

func main() {
	client := grpc_client.NewGRPCClient()
	server := handler.NewServer(":8080", client)
	fmt.Printf("Starting server at %s\n", server.Address)
	if err := server.Start(); err == nil {
		fmt.Println("Error")
	}
}
