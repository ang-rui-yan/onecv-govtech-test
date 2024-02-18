package main

import (
	"log"
	"os"
	"studentadmin/api"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	router := gin.New()

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
    handler := api.NewTeacherHandler(service)

	api := router.Group("/api")
	
	api.POST("/register", handler.RegisterHandler)
	api.GET("/commonstudents", handler.GetCommonStudentsHandler)
	api.POST("/suspend", handler.SuspendHandler)
	api.POST("/retrievefornotifications", handler.RetrieveForNotificationsHandler)

	router.Run(":8080")
}