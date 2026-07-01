package security

import (
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/microcosm-cc/bluemonday"
	"golang.org/x/time/rate"
	"github.com/gin-contrib/cors"
	"github.com/gin-contrib/secure"
)


// =========================
// GLOBAL TOOLS
// =========================

var (
	validate = validator.New()
	sanitizer = bluemonday.UGCPolicy()
)


// =========================
// RATE LIMIT (IN-MEMORY)
// =========================

type client struct {
	limiter  *rate.Limiter
	lastSeen time.Time
}



var (
	clients = make(map[string]*client)
	mu      sync.Mutex

	// cleanup config
	CLEANUP_INTERVAL = 10 * time.Minute
	MAX_IDLE         = 15 * time.Minute
)

func getClient(ip string) *rate.Limiter {
	mu.Lock()
	defer mu.Unlock()

	if c, exists := clients[ip]; exists {
		c.lastSeen = time.Now()
		return c.limiter
	}

	limiter := rate.NewLimiter(2, 5)

	clients[ip] = &client{
		limiter:  limiter,
		lastSeen: time.Now(),
	}

	return limiter
}


// =========================
// SECURITY MIDDLEWARE STACK
// =========================

func SecurityStack() []gin.HandlerFunc {
	return []gin.HandlerFunc{
		SecurityHeaders(),
		Cors(),
		RateLimiter(),
	}
}


// =========================
// SECURITY HEADERS
// =========================

func SecurityHeaders() gin.HandlerFunc {
	return secure.New(secure.Config{
		FrameDeny:             true,
		ContentTypeNosniff:    true,
		BrowserXssFilter:      true,
		ReferrerPolicy:        "strict-origin-when-cross-origin",
	})
}


// =========================
// CORS
// =========================

func Cors() gin.HandlerFunc {
	return cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:3000"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE"},
		AllowHeaders:     []string{"Authorization", "Content-Type"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	})
}


// =========================
// RATE LIMIT MIDDLEWARE
// =========================

func RateLimiter() gin.HandlerFunc {
	return func(c *gin.Context) {

		ip := c.ClientIP()
		limiter := getClient(ip)

		if !limiter.Allow() {
			c.AbortWithStatusJSON(http.StatusTooManyRequests, gin.H{
				"error": "too many requests",
			})
			return
		}

		c.Next()
	}
}


// =========================
// INPUT SECURITY (CLEAN + SAFE)
// =========================

// sanitize semua string input
func Sanitize(input string) string {
	input = strings.TrimSpace(input)
	return sanitizer.Sanitize(input)
}


// validate email
func IsValidEmail(email string) bool {
	err := validate.Var(email, "required,email")
	return err == nil
}


// basic anti injection check
func IsSafe(input string) bool {
	lower := strings.ToLower(input)

	blocked := []string{
		"<script", "javascript:", "onerror=", "onload=",
		"select ", "insert ", "drop ", "--", ";",
	}

	for _, b := range blocked {
		if strings.Contains(lower, b) {
			return false
		}
	}

	return true
}


// =========================
// VALIDATION WRAPPER (OPTIONAL)
// =========================

func ValidateStruct(s interface{}) error {
	return validate.Struct(s)
}


func StartCleanupRoutine() {
	go func() {
		for {
			time.Sleep(CLEANUP_INTERVAL)

			mu.Lock()

			now := time.Now()
			for ip, c := range clients {
				if now.Sub(c.lastSeen) > MAX_IDLE {
					delete(clients, ip)
				}
			}

			mu.Unlock()
		}
	}()
}