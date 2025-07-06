package main

import (
	"context"
	"golang-case/internal/adapter/http"
	"log"
	"os"
	"time"

	"github.com/gofiber/fiber/v2"
	mg "go.mongodb.org/mongo-driver/mongo" // Resmi sürücüye mg alias’ı
	"go.mongodb.org/mongo-driver/mongo/options"

	repo "golang-case/internal/adapter/repository" // Senin adapter paketine repo alias’ı
	"golang-case/internal/core/service"
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
	svc := service.NewDriverLocationService(repository)

	http.RegisterDriverRoutes(app, svc)

	log.Println("Listening on :3000")
	log.Fatal(app.Listen(":3000"))
}
