package handler

import (
	"context"
	"crypto/rand"
	"encoding/hex"

	"net/http"
	"time"

	"oauth-golang/internal/database"
	"oauth-golang/internal/mail"
	"oauth-golang/internal/security"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

type RegisterRequest struct {
	Username string `json:"username" binding:"required"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
}

// =========================
// TOKEN GENERATOR (SAFE)
// =========================
func generateToken() (string, error) {
	b := make([]byte, 32)

	_, err := rand.Read(b)
	if err != nil {
		return "", err
	}

	return hex.EncodeToString(b), nil
}

// =========================
// REGISTER
// =========================
func Register(c *gin.Context) {
	var req RegisterRequest

	// 1. BIND
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 2. SANITIZE
	req.Username = security.Sanitize(req.Username)
	req.Email = security.Sanitize(req.Email)

	// 3. VALIDATE
	if !security.IsValidEmail(req.Email) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid email"})
		return
	}

	if !security.IsSafe(req.Username) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid username"})
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// 4. CHECK EMAIL EXISTS
	var exists bool
	err := database.DB.QueryRow(ctx, `
		SELECT EXISTS(
			SELECT 1
			FROM users
			WHERE email = $1
		)
	`, req.Email).Scan(&exists)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "db error"})
		return
	}

	if exists {
		c.JSON(http.StatusConflict, gin.H{"error": "email already exists"})
		return
	}

	// 5. HASH PASSWORD
	hash, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "hash error"})
		return
	}

	// 6. START TRANSACTION
	tx, err := database.DB.Begin(ctx)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "tx error"})
		return
	}
	defer tx.Rollback(ctx)

	// 7. INSERT USER
	var userID string

	err = tx.QueryRow(ctx, `
		INSERT INTO users (username, email, password_hash, email_verified)
		VALUES ($1, $2, $3, false)
		RETURNING id
	`, req.Username, req.Email, string(hash)).Scan(&userID)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "insert user failed"})
		return
	}

	// 8. GENERATE TOKEN
	token, err := generateToken()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "token error"})
		return
	}

	// 9. SAVE TOKEN
	_, err = tx.Exec(ctx, `
		INSERT INTO email_verifications (user_id, token, expires_at)
		VALUES ($1, $2, NOW() + INTERVAL '15 minutes')
	`, userID, token)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "token insert failed",
			"details": err.Error(),
		})
		return
	}

	// 10. COMMIT TRANSACTION
	if err := tx.Commit(ctx); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "commit failed"})
		return
	}

	// 11. SEND VERIFICATION EMAIL
	err = mail.SendVerificationEmail(req.Email, token)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "failed send verification email",
		})
		return
	}

	// 12. SUCCESS
	c.JSON(http.StatusCreated, gin.H{
		"message": "Registration successful. Please check your email.",
	})
}

// =========================
// VERIFY EMAIL (OUTSIDE REGISTER)
// =========================
func VerifyEmail(c *gin.Context) {
	token := c.Query("token")

	var userID string

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
defer cancel()

err := database.DB.QueryRow(ctx, `
    SELECT user_id
    FROM email_verifications
    WHERE token=$1
    AND expires_at > NOW()
`, token).Scan(&userID)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid or expired token"})
		return
	}

	// update user
	_, err = database.DB.Exec(c, `
		UPDATE users SET email_verified=true WHERE id=$1
	`, userID)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "update failed"})
		return
	}

	// delete token
	_, _ = database.DB.Exec(ctx, `
    DELETE FROM email_verifications
    WHERE token=$1
`, token)
}