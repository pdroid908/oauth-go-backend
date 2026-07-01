package handler

import (
	
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	
	"net/http"
	"os"
	"github.com/joho/godotenv"
	"time"

	"oauth-golang/internal/database"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

var googleOauthConfig *oauth2.Config

func init() {
    godotenv.Load()

    googleOauthConfig = &oauth2.Config{
        ClientID:     os.Getenv("GOOGLE_CLIENT_ID"),
        ClientSecret: os.Getenv("GOOGLE_CLIENT_SECRET"),
        RedirectURL:  "https://pasdaoji-backend-oauth.hf.space/auth/google/callback",
        Scopes: []string{
            "https://www.googleapis.com/auth/userinfo.email",
            "https://www.googleapis.com/auth/userinfo.profile",
        },
        Endpoint: google.Endpoint,
    }
}

// generateState membuat string acak untuk keamanan CSRF
func generateState(c *gin.Context) string {
	b := make([]byte, 16)
	rand.Read(b)
	state := base64.URLEncoding.EncodeToString(b)
	c.SetCookie("oauthstate", state, 3600, "/", "", false, true)
	return state
}

func GoogleLogin(c *gin.Context) {
	state := generateState(c)
	url := googleOauthConfig.AuthCodeURL(state)
	c.Redirect(http.StatusTemporaryRedirect, url)
}

func GoogleCallback(c *gin.Context) {
	// 1. Verifikasi State (CSRF Protection)
	state, err := c.Cookie("oauthstate")
	if err != nil || c.Query("state") != state {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid state"})
		return
	}

	// 2. Exchange Code to Token
	code := c.Query("code")
	token, err := googleOauthConfig.Exchange(c.Request.Context(), code)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to exchange token"})
		return
	}

	// 3. Ambil data user dari Google
	client := googleOauthConfig.Client(c.Request.Context(), token)
	resp, err := client.Get("https://www.googleapis.com/oauth2/v2/userinfo")
	if err != nil || resp.StatusCode != http.StatusOK {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get user info"})
		return
	}
	defer resp.Body.Close()

	var userInfo struct {
		Email string `json:"email"`
		Name  string `json:"name"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&userInfo); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to decode user info"})
		return
	}

	// 4. Upsert User ke Database
	_, err = database.DB.Exec(c.Request.Context(), `
		INSERT INTO users (username, email, email_verified, password_hash)
		VALUES ($1, $2, true, 'google_user')
		ON CONFLICT (email) DO UPDATE SET email_verified = true
	`, userInfo.Name, userInfo.Email)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
		return
	}

	// 5. Generate JWT
	secret := []byte(os.Getenv("JWT_SECRET"))
	claims := jwt.MapClaims{
		"username": userInfo.Name,
		"email":    userInfo.Email,
		"exp":      time.Now().Add(time.Hour * 1).Unix(),
	}
	tokenString, _ := jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString(secret)

	// 6. Set Cookie dan Response
	c.SetCookie("token", tokenString, 3600, "/", "", false, true)
	c.JSON(http.StatusOK, gin.H{"message": "Login successful", "token": tokenString})
}