package handler

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"net/http"
	"time"
	"oauth-golang/internal/database"
	"oauth-golang/internal/mail" 
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)


// ForgotPassword handle request awal user lupa password
func ForgotPassword(c *gin.Context) {
	var req struct {
		Email string `json:"email" binding:"required,email"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid email format"})
		return
	}

	var userID string
	err := database.DB.QueryRow(c, "SELECT id FROM users WHERE email = $1", req.Email).Scan(&userID)

	if err != nil {
		// Tetap berikan response sukses agar tidak membocorkan email yang terdaftar
		c.JSON(http.StatusOK, gin.H{"message": "If the email exists, a reset link has been sent."})
		return
	}

	// Buat Token
	b := make([]byte, 32)
	rand.Read(b)
	token := hex.EncodeToString(b)
	expiresAt := time.Now().Add(time.Hour)

	// Hapus token lama jika ada agar tidak duplikat
	_, _ = database.DB.Exec(c, "DELETE FROM password_resets WHERE user_id = $1", userID)

	// Insert token baru (4 kolom: id, user_id, token, expires_at)
	_, err = database.DB.Exec(
		c,
		`INSERT INTO password_resets (id, user_id, token, expires_at) 
		 VALUES (gen_random_uuid(), $1, $2, $3)`,
		userID,
		token,
		expiresAt,
	)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save reset token"})
		return
	}

	// Kirim email
	err = mail.SendResetPasswordEmail(req.Email, token)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to send email"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Reset link sent"})
    c.Redirect(http.StatusFound, "https://netizencom.pages.dev/tampilan/reset-password")
}


// ResetPassword handle form submit password baru
func ResetPassword(c *gin.Context) {
    var req struct {
        Token       string `json:"token" binding:"required"`
        NewPassword string `json:"new_password" binding:"required,min=6"`
    }
    
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
        return
    }

    // 1. Cari data di password_resets
    // PERBAIKAN: userID diubah dari int menjadi string (karena UUID di DB adalah text/uuid)
    var userID string 
    var expiresAt time.Time
    err := database.DB.QueryRow(c, 
        "SELECT user_id, expires_at FROM password_resets WHERE token = $1", req.Token).Scan(&userID, &expiresAt)
    
    if err != nil {
        // Log error untuk debug di terminal
        fmt.Println("Query token error:", err)
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid or expired token"})
        return
    }

    // 2. Cek expired
    if time.Now().After(expiresAt) {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Token expired"})
        return
    }

    // 3. Hash password baru
    hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.NewPassword), bcrypt.DefaultCost)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Hashing failed"})
        return
    }

    // 4. Transaksi: Update password dan Hapus token
    tx, err := database.DB.Begin(c)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Transaction failed"})
        return
    }
    defer tx.Rollback(c)

    // Update password
    _, err = tx.Exec(c, "UPDATE users SET password_hash = $1 WHERE id = $2", string(hashedPassword), userID)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Update password failed"})
        return
    }

    // Hapus token setelah digunakan
    _, err = tx.Exec(c, "DELETE FROM password_resets WHERE token = $1", req.Token)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Cleanup failed"})
        return
    }

    // Komit transaksi
    err = tx.Commit(c)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Commit failed"})
        return
    }

    c.JSON(http.StatusOK, gin.H{"message": "Password updated successfully"})
    c.Redirect(http.StatusFound, "https://netizencom.pages.dev/tampilan/login")
    
}