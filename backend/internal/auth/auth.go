package auth

import (
	"context"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	ID           uuid.UUID
	Username     string
	PasswordHash string
}

// LoginRequest структура входных данных для авторизации
type LoginRequest struct {
	Username string `json:"username" example:"admin"`
	Password string `json:"password" example:"123456"`
}

// LoginResponse структура ответа с токеном
type LoginResponse struct {
	Token string `json:"token" example:"eyJhbGciOi..."`
}

// Login godoc
// @Summary Авторизация
// @Description Авторизация пользователя и получение JWT
// @Tags auth
// @Accept json
// @Produce json
// @Param credentials body auth.LoginRequest true "Учетные данные"
// @Success 200 {object} auth.LoginResponse
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Router /api/auth/login [post]
func Login(pool *pgxpool.Pool) gin.HandlerFunc {
	return func(c *gin.Context) {
		var input LoginRequest
		if err := c.ShouldBindJSON(&input); err != nil {
			log.Printf("Failed to parse login input: %v", err)
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input format"})
			return
		}

		var user User
		row := pool.QueryRow(context.Background(), "SELECT id, password_hash FROM users WHERE username=$1", input.Username)
		if err := row.Scan(&user.ID, &user.PasswordHash); err != nil {
			log.Printf("User not found or query error for username=%s: %v", input.Username, err)
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
			return
		}

		if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(input.Password)); err != nil {
			log.Printf("Invalid password or bcrypt error: %v", err)
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
			return
		}

		jwtSecret := os.Getenv("JWT_SECRET")
		if jwtSecret == "" {
			log.Printf("JWT_SECRET is not set in environment variables")
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Configuration error"})
			return
		}

		token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
			"sub": user.ID.String(),
			"exp": time.Now().Add(72 * time.Hour).Unix(),
		})

		tokenString, err := token.SignedString([]byte(jwtSecret))
		if err != nil {
			log.Printf("Failed to sign JWT: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not generate token"})
			return
		}

		c.JSON(http.StatusOK, LoginResponse{Token: tokenString})
	}
}
