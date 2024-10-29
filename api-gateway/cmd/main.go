package main

import (
	"fmt"

	"github.com/lunarKettle/task-management-platform/api-gateway/internal/grpc_client"
	"github.com/lunarKettle/task-management-platform/api-gateway/internal/handler"
)

func main() {
	client := grpc_client.NewGRPCClient()
	server := handler.NewServer(":8080", client)
	fmt.Printf("Starting server at %s\n", server.Address)
	if err := server.Start(); err == nil {
		fmt.Println("Error")
	}
}
