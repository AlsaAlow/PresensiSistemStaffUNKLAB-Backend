package main

import (
	"net/http"
	"time"

	"sipres/config"
	"sipres/internal/handler"
	"sipres/internal/middleware"
	"sipres/internal/repository"
	"sipres/internal/service"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

// 🔐 HARUS SAMA dengan middleware
var jwtSecret = []byte("secret_key")

func main() {
	// connect database
	config.ConnectDB()

	// setup clean architecture
	repo := &repository.AttendanceRepository{DB: config.DB}
	svc := &service.AttendanceService{Repo: repo}
	h := &handler.AttendanceHandler{Service: svc}

	r := gin.Default()

	// =========================
	// TEST
	// =========================
	r.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "SIPRES API Running",
		})
	})

	// =========================
	// REGISTER
	// =========================
	r.POST("/register", func(c *gin.Context) {
		var input struct {
			Name     string `json:"name"`
			Email    string `json:"email"`
			Password string `json:"password"`
		}

		if err := c.ShouldBindJSON(&input); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid input"})
			return
		}

		if input.Email == "" || input.Password == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "email/password required"})
			return
		}

		hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)

		result, err := config.DB.Exec(
			"INSERT INTO users (name, email, password, role) VALUES (?, ?, ?, ?)",
			input.Name, input.Email, string(hashedPassword), "user",
		)

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		id, _ := result.LastInsertId()

		c.JSON(http.StatusOK, gin.H{
			"message": "register success",
			"id":      id,
		})
	})

	// =========================
	// LOGIN
	// =========================
	r.POST("/login", func(c *gin.Context) {
		var input struct {
			Email    string `json:"email"`
			Password string `json:"password"`
		}

		c.ShouldBindJSON(&input)

		var user struct {
			ID       int
			Name     string
			Email    string
			Password string
			Role     string
		}

		err := config.DB.QueryRow(
			"SELECT id, name, email, password, role FROM users WHERE email = ?",
			input.Email,
		).Scan(&user.ID, &user.Name, &user.Email, &user.Password, &user.Role)

		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "user not found"})
			return
		}

		if bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(input.Password)) != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "wrong password"})
			return
		}

		token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
			"user_id": user.ID,
			"email":   user.Email,
			"role":    user.Role,
			"exp":     time.Now().Add(24 * time.Hour).Unix(),
		})

		tokenString, _ := token.SignedString(jwtSecret)

		c.JSON(http.StatusOK, gin.H{
			"message": "login success",
			"token":   tokenString,
		})
	})

	// =========================
	// PROFILE
	// =========================
	r.GET("/profile", middleware.AuthMiddleware(), func(c *gin.Context) {
		userID := int(c.GetFloat64("user_id"))

		var user struct {
			ID    int
			Name  string
			Email string
			Role  string
		}

		err := config.DB.QueryRow(
			"SELECT id, name, email, role FROM users WHERE id = ?",
			userID,
		).Scan(&user.ID, &user.Name, &user.Email, &user.Role)

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "user not found"})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"message": "profile fetched",
			"user":    user,
		})
	})

	// =========================
	// ATTENDANCE (CLEAN)
	// =========================
	r.POST("/checkin", middleware.AuthMiddleware(), h.Checkin)
	r.POST("/checkout", middleware.AuthMiddleware(), h.Checkout)
	r.GET("/attendance", middleware.AuthMiddleware(), h.GetAttendance)

	// run server
	r.Run(":8080")
}
