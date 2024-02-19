package main

import (
	"log"
	"os"
	"studentadmin/api"
)

func main() {
	host := os.Getenv("HOST")
	port := os.Getenv("PORT")
	user := os.Getenv("USER")
	dbname := os.Getenv("DB_NAME")
	pass := os.Getenv("PASSWORD")

	dbPool, _, err := api.NewDBPool(api.DatabaseConfig{
		Username: user,
		Password: pass,
		Hostname: host,
		Port:     port,
		DBName:   dbname,
	})

	if err != nil {
		log.Fatalf("unexpected error while tried to connect to database: %v\n", err)
	}

	defer dbPool.Close()

	DB := api.NewDatabase(dbPool)
	service := api.NewTeacherService(DB)
	router := api.SetupRouter(service)

	router.Run("0.0.0.0:8080")
}
