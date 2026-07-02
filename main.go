package main

import (
	"log"

	"oauth-golang/internal/database"
	"oauth-golang/internal/handler"
	"oauth-golang/internal/login"
	"os"

	"oauth-golang/internal/security"

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

	r.Use(security.Cors())

	// LOGIN REGISTER RESET PASS


	r.POST("/register", handler.Register)
	r.GET("/verify-email", handler.VerifyEmail)
	r.POST("/login", login.Login)
	r.GET("/auth/google", handler.GoogleLogin)
	r.GET("/auth/google/callback", handler.GoogleCallback)
	r.POST("/forgot-password", handler.ForgotPassword)
    r.POST("/reset-password", handler.ResetPassword)
	
	

	r.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status": "ok",
		})
	})

port := os.Getenv("PORT")
    if port == "" {
        port = "8080" // Default untuk pengembangan lokal
    }

    log.Printf("🚀 Server running on port %s", port)
    
    // Gunakan port tersebut untuk menjalankan server
    r.Run(":" + port)
}