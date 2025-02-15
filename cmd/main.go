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

	"github.com/alastor-4/sylcot-go-gin-backend/controllers"
	"github.com/alastor-4/sylcot-go-gin-backend/database"
	"github.com/alastor-4/sylcot-go-gin-backend/models"

	"github.com/alastor-4/sylcot-go-gin-backend/middleware"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"gorm.io/gorm"
)

type Repository struct {
	DB *gorm.DB
}

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
		panic("Failed to seed categories")
	}

	ac := &controllers.AuthController{DB: db}
	tc := &controllers.TaskController{DB: db}

	router := setupRouter(ac, tc)
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	router.Run(":" + port)
}

func setupRouter(ac *controllers.AuthController, tc *controllers.TaskController) *gin.Engine {

	router := gin.Default()
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	auth := router.Group("/auth")
	{
		auth.POST("/register", ac.Register)
		auth.POST("/login", ac.Login)
		auth.GET("/verify-email/", ac.VerifyEmail)
		auth.GET("/refresh", ac.Refresh)
	}

	api := router.Group("/api").Use(middleware.AuthMiddleware())
	{
		api.GET("/tasks", tc.GetTasks)
		api.POST("/tasks", tc.CreateTask)
		api.PUT("/tasks/:id", tc.UpdateTask)
		api.DELETE("/tasks/:id", tc.DeleteTask)
		api.PATCH("/tasks/:id/complete", tc.ToggleTask)
	}

	return router
}
