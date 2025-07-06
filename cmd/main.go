package main

import (
	"context"
	"github.com/envercigal/golang/internal/adapter/http"
	repo "github.com/envercigal/golang/internal/adapter/repository"
	"github.com/envercigal/golang/internal/core/service"
	circuitbreaker "github.com/envercigal/golang/pkg"
	"log"
	"os"
	"time"

	"github.com/gofiber/fiber/v2"
	mg "go.mongodb.org/mongo-driver/mongo" // Resmi sürücüye mg alias’ı
	"go.mongodb.org/mongo-driver/mongo/options"
)

func main() {
	app := fiber.New()

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	mongoURI := os.Getenv("MONGO_URI")
	if mongoURI == "" {
		mongoURI = "mongodb://admin:password@localhost:27017"
	}

	client, err := mg.Connect(ctx, options.Client().ApplyURI(mongoURI))
	if err != nil {
		log.Fatal(err)
	}

	repository := repo.NewDriverLocationRepo(client.Database("mydb"))
	circuitBreaker := circuitbreaker.New(5, 1)
	svc := service.NewDriverLocationService(repository, circuitBreaker)

	http.RegisterDriverRoutes(app, svc)

	log.Println("Listening on :3000")
	log.Fatal(app.Listen(":3000"))
}
