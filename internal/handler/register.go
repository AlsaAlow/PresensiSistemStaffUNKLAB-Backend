package handler

import (
	"net/http"

	"sipres/config"

	"github.com/gin-gonic/gin"
)

type RegisterInput struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
	Role     string `json:"role"`
}

func Register(c *gin.Context) {
	var input RegisterInput

	// bind JSON
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid input",
		})
		return
	}

	// insert ke database
	result, err := config.DB.Exec(
		"INSERT INTO users (name, email, password, role) VALUES (?, ?, ?, ?)",
		input.Name, input.Email, input.Password, input.Role,
	)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(), // 🔥 tampilkan error asli
		})
		return
	}

	id, _ := result.LastInsertId()

	c.JSON(http.StatusOK, gin.H{
		"message": "register success",
		"id":      id,
	})
}
