package api

import (
	"os"

	"github.com/A4GOD-AMHG/sylcot-go-gin-backend/internal/handlers"
	"github.com/A4GOD-AMHG/sylcot-go-gin-backend/internal/repositories"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type App struct {
	router *gin.Engine
	db     *gorm.DB
}

func NewApp(db *gorm.DB) *App {
	app := &App{
		router: gin.Default(),
		db:     db,
	}

	app.initializeRoutes()
	return app
}

func (a *App) initializeRoutes() {
	authRepo := repositories.NewAuthRepository(a.db)
	taskRepo := repositories.NewTaskRepository(a.db)
	categoryRepo := repositories.NewCategoryRepository(a.db)

	authHandler := &handlers.AuthHandler{Repo: authRepo}
	taskHandler := handlers.NewTaskHandler(taskRepo)
	categoryHandler := handlers.NewCategoryHandler(categoryRepo)

	SetupRoutes(a.router, authHandler, taskHandler, categoryHandler)
}

func (a *App) Run() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	a.router.Run(":" + port)
}
