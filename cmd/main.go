package main

import (
	"log"
	"os"
	"studentadmin/api"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

func main() {
	host := os.Getenv("HOST")
	port := os.Getenv("PORT")
	user := os.Getenv("USER")
	dbname := os.Getenv("DB_NAME")
	pass := os.Getenv("PASSWORD")

	maxRetries := 5
	retryDelay := 5 * time.Second

	var dbPool *pgxpool.Pool

	// Attempt to connect to the database with retries.
	for i := 0; i < maxRetries; i++ {
		var err error
		dbPool, _, err = api.NewDBPool(api.DatabaseConfig{
			Username: user,
			Password: pass,
			Hostname: host,
			Port:     port,
			DBName:   dbname,
		})

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

	DB := api.NewDatabase(dbPool)
	service := api.NewTeacherService(DB)
	router := api.SetupRouter(service)

	router.Run("0.0.0.0:8080")
}
