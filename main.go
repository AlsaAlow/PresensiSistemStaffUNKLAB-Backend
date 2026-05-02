package main

import (
	"net/http"
	"time"

	"sipres/config"
	"sipres/internal/middleware"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

// 🔐 HARUS SAMA dengan middleware
var jwtSecret = []byte("secret_key")

func main() {
	// connect database
	config.ConnectDB()

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

		// hash password
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed hashing"})
			return
		}

		// insert ke database
		result, err := config.DB.Exec(
			"INSERT INTO users (name, email, password, role) VALUES (?, ?, ?, ?)",
			input.Name, input.Email, string(hashedPassword), "user",
		)

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": err.Error(),
			})
			return
		}

		id, _ := result.LastInsertId()

		c.JSON(http.StatusOK, gin.H{
			"message": "register success",
			"id":      id,
		})
	})

	// =========================
	// LOGIN + JWT
	// =========================
	r.POST("/login", func(c *gin.Context) {
		var input struct {
			Email    string `json:"email"`
			Password string `json:"password"`
		}

		if err := c.ShouldBindJSON(&input); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid input"})
			return
		}

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

		// compare password
		err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(input.Password))
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "wrong password"})
			return
		}

		// =========================
		// GENERATE JWT
		// =========================
		token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
			"user_id": user.ID,
			"email":   user.Email,
			"role":    user.Role,
			"exp":     time.Now().Add(time.Hour * 24).Unix(),
		})

		tokenString, err := token.SignedString(jwtSecret)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed generate token"})
			return
		}

		// =========================
		// RESPONSE
		// =========================
		c.JSON(http.StatusOK, gin.H{
			"message": "login success",
			"token":   tokenString,
			"user": gin.H{
				"id":    user.ID,
				"name":  user.Name,
				"email": user.Email,
				"role":  user.Role,
			},
		})
	})

	// =========================
	// PROTECTED ROUTE
	// =========================
	r.GET("/profile", middleware.AuthMiddleware(), func(c *gin.Context) {

		// ambil user_id dari middleware
		userID := c.GetFloat64("user_id")

		var user struct {
			ID    int
			Name  string
			Email string
			Role  string
		}

		// ambil dari database
		err := config.DB.QueryRow(
			"SELECT id, name, email, role FROM users WHERE id = ?",
			int(userID),
		).Scan(&user.ID, &user.Name, &user.Email, &user.Role)

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "user not found"})
			return
		}

		// response
		c.JSON(http.StatusOK, gin.H{
			"message": "profile fetched",
			"user":    user,
		})
	})

	// run server
	r.Run(":8080")
}
