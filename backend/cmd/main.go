package main

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	_ "license-service/docs"
	"license-service/internal/auth"
	"license-service/internal/db"
	"license-service/internal/license"
	"license-service/internal/middleware"
	"log"
	"time"
)

// @title License Service API
// @version 1.0
// @description API для управления лицензиями
// @host localhost:8080
// @BasePath /

func main() {
	_ = godotenv.Load()
	pool := db.InitDB()
	defer pool.Close()

	r := gin.Default()

	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"}, // или конкретные домены: []string{"https://example.com"}
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	r.GET("/api/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	r.POST("/api/auth/login", auth.Login(pool))

	authGroup := r.Group("/api")
	authGroup.Use(middleware.JWTMiddleware())
	{
		authGroup.POST("/license/generate", license.Generate(pool))
		authGroup.POST("/license/verify", license.Verify(pool))
	}

	adminGroup := r.Group("/api/admin")
	adminGroup.Use(middleware.JWTMiddleware())
	{
		adminGroup.GET("/licenses", license.ListLicenses(pool))

		adminGroup.GET("/licenses/:id", license.GetLicense(pool))

		adminGroup.PUT("/licenses/:id", license.UpdateLicense(pool))

		adminGroup.DELETE("/licenses/:id", license.DeleteLicense(pool))

		adminGroup.GET("/licenses/:id/download", license.DownloadLicense(pool))
	}

	log.Fatal(r.Run(":8080"))
}
