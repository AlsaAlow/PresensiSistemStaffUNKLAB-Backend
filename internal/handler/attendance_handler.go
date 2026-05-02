package handler

import (
	"net/http"

	"sipres/internal/service"

	"github.com/gin-gonic/gin"
)

type AttendanceHandler struct {
	Service *service.AttendanceService
}

// =========================
// CHECK IN
// =========================
func (h *AttendanceHandler) Checkin(c *gin.Context) {
	userID := int(c.GetFloat64("user_id"))

	msg, err := h.Service.Checkin(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if msg == "masih belum check-out" {
		c.JSON(http.StatusBadRequest, gin.H{"error": msg})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": msg,
	})
}

// =========================
// CHECK OUT
// =========================
func (h *AttendanceHandler) Checkout(c *gin.Context) {
	userID := int(c.GetFloat64("user_id"))

	msg, err := h.Service.Checkout(userID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": msg,
	})
}

// =========================
// GET ATTENDANCE
// =========================
func (h *AttendanceHandler) GetAttendance(c *gin.Context) {
	userID := int(c.GetFloat64("user_id"))

	data, err := h.Service.GetAttendance(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": data,
	})
}
