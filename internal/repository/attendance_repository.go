package repository

import (
	"database/sql"
)

type AttendanceRepository struct {
	DB *sql.DB
}

// =========================
// CEK CHECKIN AKTIF
// =========================
func (r *AttendanceRepository) CheckActiveCheckin(userID int) (bool, error) {
	var id int

	err := r.DB.QueryRow(
		`SELECT id FROM attendances 
		 WHERE user_id = ? 
		 AND DATE(check_in) = CURDATE()
		 AND check_out IS NULL`,
		userID,
	).Scan(&id)

	if err == sql.ErrNoRows {
		return false, nil
	}
	if err != nil {
		return false, err
	}

	return true, nil
}

// =========================
// INSERT CHECKIN
// =========================
func (r *AttendanceRepository) InsertCheckin(userID int) error {
	_, err := r.DB.Exec(
		"INSERT INTO attendances (user_id, check_in) VALUES (?, NOW())",
		userID,
	)
	return err
}

// =========================
// UPDATE CHECKOUT
// =========================
func (r *AttendanceRepository) UpdateCheckout(userID int) (int64, error) {
	result, err := r.DB.Exec(
		`UPDATE attendances 
		 SET check_out = NOW() 
		 WHERE user_id = ? 
		 AND DATE(check_in) = CURDATE()
		 AND check_out IS NULL`,
		userID,
	)

	if err != nil {
		return 0, err
	}

	return result.RowsAffected()
}

// =========================
// GET ATTENDANCE
// =========================
func (r *AttendanceRepository) GetAttendance(userID int) ([]map[string]interface{}, error) {
	rows, err := r.DB.Query(
		"SELECT id, check_in, check_out FROM attendances WHERE user_id = ?",
		userID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var data []map[string]interface{}

	for rows.Next() {
		var id int
		var checkIn, checkOut sql.NullTime

		err := rows.Scan(&id, &checkIn, &checkOut)
		if err != nil {
			return nil, err
		}

		var checkInVal interface{}
		var checkOutVal interface{}

		if checkIn.Valid {
			checkInVal = checkIn.Time.Format("2006-01-02 15:04:05")
		}

		if checkOut.Valid {
			checkOutVal = checkOut.Time.Format("2006-01-02 15:04:05")
		}

		data = append(data, map[string]interface{}{
			"id":        id,
			"check_in":  checkInVal,
			"check_out": checkOutVal,
		})
	}

	return data, nil
}
