package main

import (
	"log"

	"oauth-golang/internal/database"
	"oauth-golang/internal/handler"
	"oauth-golang/internal/login"

	"oauth-golang/internal/security"
	"github.com/gin-contrib/cors"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	// Load .env
	if err := godotenv.Load(); err != nil {
		log.Println(".env not found, using system environment")
	}

	// Connect PostgreSQL
	database.Connect()
	defer database.Close()

	// Gin
	r := gin.Default()

	r.Use(security.SecurityStack()...)
	security.StartCleanupRoutine()

	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:3000"}, // URL Frontend Anda
		AllowMethods:     []string{"POST", "GET", "OPTIONS", "DELETE"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true, // WAJIB: agar cookie/session bisa dikirim
		
	}))

	// LOGIN REGISTER RESET PASS


	r.POST("/register", handler.Register)
	r.GET("/verify-email", handler.VerifyEmail)
	r.POST("/login", login.Login)
	r.GET("/auth/google", handler.GoogleLogin)
	r.GET("/auth/google/callback", handler.GoogleCallback)
	r.POST("/forgot-password", handler.ForgotPassword)
    r.POST("/reset-password", handler.ResetPassword)
	
	

	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status": "ok",
		})
	})

	log.Println("🚀 Server running on :8080")
	r.Run(":8080")
}