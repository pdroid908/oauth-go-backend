package login

import (
	"context"
	"net/http"
	"os"
	"time"

	"oauth-golang/internal/database"
	"oauth-golang/internal/security"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

type LoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type Claims struct {
	Username string `json:"username"`
	jwt.RegisteredClaims
}

func Login(c *gin.Context) {
	var req LoginRequest

	// 1. BIND & VALIDATE
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request format"})
		return
	}

	// 2. SANITIZE
	req.Username = security.Sanitize(req.Username)

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	// 3. FETCH USER (Cari berdasarkan Username)
	var passwordHash string
	var emailVerified bool
	err := database.DB.QueryRow(ctx, `
		SELECT password_hash, email_verified 
		FROM users 
		WHERE username = $1
	`, req.Username).Scan(&passwordHash, &emailVerified)

	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid credentials"})
		return
	}

	// 4. CHECK VERIFIED
	if !emailVerified {
		c.JSON(http.StatusForbidden, gin.H{"error": "please verify your email first"})
		return
	}

	// 5. CHECK PASSWORD
	if err := bcrypt.CompareHashAndPassword([]byte(passwordHash), []byte(req.Password)); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid credentials"})
		return
	}

	// 6. GENERATE JWT
	secret := []byte(os.Getenv("JWT_SECRET"))
	if len(secret) == 0 {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "server configuration error"})
		return
	}

	expirationTime := time.Now().Add(1 * time.Hour)
	claims := &Claims{
		Username: req.Username,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(secret)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "could not generate token"})
		return
	}

	// 7. SET COOKIE
	c.SetCookie(
		"token",
		tokenString,
		3600, // 1 hour
		"/",
		"",
		true, // Set ke true jika sudah pakai HTTPS
		true,  // HttpOnly
	)

	c.JSON(http.StatusOK, gin.H{
    "success": true, 
    "message": "Login berhasil",
})

	redirectURL := "https://netizencom.pages.dev/api/callback?token=" + tokenString
	c.Redirect(http.StatusFound, redirectURL)
}