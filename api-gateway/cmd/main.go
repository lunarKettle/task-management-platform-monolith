package main

import (
	"cmd/main.go/internal/handler"
)

func main() {
	server := handler.NewServer(":8080")
	server.Start()
}
