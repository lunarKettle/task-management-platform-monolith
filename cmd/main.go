package main

import (
	"database/sql"
	"fmt"

	userInfrastructure "github.com/lunarKettle/task-management-platform-monolith/internal/user/infrastructure"
	userTransport "github.com/lunarKettle/task-management-platform-monolith/internal/user/transport"
	userUsecases "github.com/lunarKettle/task-management-platform-monolith/internal/user/usecases"

	projectInfrastructure "github.com/lunarKettle/task-management-platform-monolith/internal/project/infrastructure"
	projectTransport "github.com/lunarKettle/task-management-platform-monolith/internal/project/transport"
	projectUsecases "github.com/lunarKettle/task-management-platform-monolith/internal/project/usecases"

	"github.com/lunarKettle/task-management-platform-monolith/internal/server"

	"github.com/joho/godotenv"

	"log"
	"os"

	_ "github.com/lib/pq"
)

const (
	ServerAddressEnv    = "SERVER_ADDRESS"
	ConnectionStringEnv = "CONNECTION_STRING"
	SecretKeyEnv        = "SECRET_KEY"
)

func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
}

func run() error {
	if err := godotenv.Load(); err != nil {
		return fmt.Errorf("failed to load .env file: %w", err)
	}

	database, err := sql.Open("postgres", os.Getenv(ConnectionStringEnv))
	if err != nil {
		return fmt.Errorf("failed to open connection: %w", err)
	}

	userRepo := userInfrastructure.NewUserRepository(database)
	jwtManager := userInfrastructure.NewJWTManager(os.Getenv(SecretKeyEnv))
	projectRepo := projectInfrastructure.NewProjectRepository(database)

	authUseCases := userUsecases.NewAuthUseCases(userRepo, jwtManager)
	projectUseCases := projectUsecases.NewProjectUseCases(projectRepo)

	authHandlers := userTransport.NewAuthHandlers(authUseCases)
	projectHandler := projectTransport.NewProjectHandlers(projectUseCases)

	server := server.NewServer(":8080")
	fmt.Printf("Starting server at %s\n", server.Address)
	if err := server.Start(authHandlers, projectHandler); err == nil {
		fmt.Println("Error")
	}
	return nil
}
