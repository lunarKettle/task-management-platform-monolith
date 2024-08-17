package main

import (
	"api-gateway/internal/grpc_client"
	"log"
	"time"
)

func main() {
	// server := handler.NewServer(":8080")
	// fmt.Printf("Starting server at %s\n", server.Address)
	// if err := server.Start(); err == nil {
	// 	fmt.Println("Ошибка")
	// }

	client := grpc_client.NewGRPCClient()
	r, _ := client.GetProject()
	log.Printf("Projects: %s", r.GetProjectDescription())
	time.Sleep(10 * time.Second)
}
