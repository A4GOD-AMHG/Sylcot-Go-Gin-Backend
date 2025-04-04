// @title Sylcot API
// @version 1.0
// @description Sylcot, that is an acronym for Simplify Your Life by Crossing Out Tasks, it is Task management API to manage your priorities, with a little more functionality and complexity, like JWT authentication

// @contact.email alexismhgarcia@gmail.com

// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name Authorization

package main

import (
	"log"
	"os"
	"time"

	"github.com/A4GOD-AMHG/sylcot-go-gin-backend/cmd/api"
	_ "github.com/A4GOD-AMHG/sylcot-go-gin-backend/docs"
	"github.com/A4GOD-AMHG/sylcot-go-gin-backend/internal/models"
	"github.com/A4GOD-AMHG/sylcot-go-gin-backend/pkg/database"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {
	config := database.Config{
		Host:     os.Getenv("DB_HOST"),
		Port:     os.Getenv("DB_PORT"),
		Password: os.Getenv("DB_PASSWORD"),
		User:     os.Getenv("DB_USER"),
		DBName:   os.Getenv("DB_NAME"),
	}

	db, err := database.NewConnection(&config)
	if err != nil {
		log.Fatal("Could not connect the database")
	}

	if err := models.MigrateAll(db); err != nil {
		log.Fatal("Migration failed: ", err)
	}

	var category models.Category
	if err := category.Setup(db); err != nil {
		log.Fatal("Failed to seed categories: ", err)
	}

	router := gin.Default()
	corsConfig := cors.Config{
		AllowOrigins:     []string{"http://localhost:5173"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization", "X-Requested-With"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}
	router.Use(cors.New(corsConfig))

	app := api.NewAppWithRouter(router, db)

	app.Run()
}
