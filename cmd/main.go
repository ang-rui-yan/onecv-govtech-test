package main

import (
	"log"
	"os"
	"studentadmin/api"
	"studentadmin/cmd/cli"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
)

func main() {
	godotenv.Load("./.env")
	dbURL := os.Getenv("POSTGRESQL_URL")

	maxRetries := 5
	retryDelay := 5 * time.Second

	var dbPool *pgxpool.Pool

	// Attempt to connect to the database with retries.
	for i := 0; i < maxRetries; i++ {
		var err error
		dbPool, _, err = api.NewDBPool(dbURL)

		if err == nil {
			break
		}

		log.Printf("Failed to connect to database: %v\n", err)
		if i < maxRetries-1 {
			log.Printf("Retrying in %v...\n", retryDelay)
			time.Sleep(retryDelay)
		} else {
			log.Fatal("Could not connect to database after retries")
		}
	}

	defer dbPool.Close()

	cli.Execute(dbURL)

	DB := api.NewDatabase(dbPool)
	service := api.NewTeacherService(DB)
	router := api.SetupRouter(service)

	router.Run("0.0.0.0:8080")
}
