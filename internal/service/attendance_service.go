package service

import (
	"errors"
	"sipres/internal/repository"
)

type AttendanceService struct {
	Repo *repository.AttendanceRepository
}

// =========================
// CHECK IN
// =========================
func (s *AttendanceService) Checkin(userID int) (string, error) {
	exists, err := s.Repo.CheckActiveCheckin(userID)
	if err != nil {
		return "", err
	}

	if exists {
		return "masih belum check-out", nil
	}

	err = s.Repo.InsertCheckin(userID)
	if err != nil {
		return "", err
	}

	return "check-in berhasil", nil
}

// =========================
// CHECK OUT
// =========================
func (s *AttendanceService) Checkout(userID int) (string, error) {
	affected, err := s.Repo.UpdateCheckout(userID)
	if err != nil {
		return "", err
	}

	if affected == 0 {
		return "", errors.New("belum check-in atau sudah check-out")
	}

	return "check-out berhasil", nil
}

// =========================
// GET ATTENDANCE
// =========================
func (s *AttendanceService) GetAttendance(userID int) ([]map[string]interface{}, error) {
	return s.Repo.GetAttendance(userID)
}
