package main

import (
	"log"
	"os"
	"studentadmin/api"

	"github.com/joho/godotenv"
)

func main() {
	var err error
	err = godotenv.Load()
	if err != nil {
		log.Println("Error has occurred on .env file. Please check.")
	}

	host := os.Getenv("HOST")
	port := os.Getenv("PORT")
	user := os.Getenv("USER")
	dbname := os.Getenv("DB_NAME")
	pass := os.Getenv("PASSWORD")

	dbPool, _, err := api.NewDBPool(api.DatabaseConfig{
        Username : user,
        Password : pass,
        Hostname : host,
        Port : port,
        DBName : dbname,
    })

    if err != nil {
        log.Fatalf("unexpected error while tried to connect to database: %v\n", err)
    }

    defer dbPool.Close()

	DB := api.NewDatabase(dbPool)
	service := api.NewTeacherService(DB)
	router := api.SetupRouter(service)

	router.Run("localhost:8080")
}