package databases

import (
	"database/sql"
	"strconv"
	"strings"
	"time"

	"github.com/RYANCOAL9999/SpinnrTechnologyInterview/gameRoomManagmentSystem/models"
)

func aizuArray(A string) []int {
	strs := strings.Split(A, ", ")
	ary := make([]int, len(strs))
	for i := range ary {
		ary[i], _ = strconv.Atoi(strs[i])
	}
	return ary
}

func ListReservation(db *sql.DB, roomID int, startDate, endDate time.Time, limit int) ([]models.ReservationRoom, error) {
	query := `
        SELECT
            r.id AS reservation_id,
            rm.id AS room_id,
            r.date_time AS reservation_date_time,
            rm.player_ids AS player_ids
        FROM reservations r
        INNER JOIN room rm ON r.room_id = rm.id
        WHERE 1=1
    `
	args := []interface{}{}

	if roomID != 0 {
		query += " AND room_id = ?"
		args = append(args, roomID)
	}

	if !startDate.IsZero() && !endDate.IsZero() {
		query += " AND date BETWEEN ? AND ?"
		args = append(args, startDate, endDate)
	}

	query += " ORDER BY reservation_id"

	if limit > 0 {
		query += " LIMIT ?"
		args = append(args, limit)
	}

	rows, err := db.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var reservations []models.ReservationRoom
	for rows.Next() {
		var r models.ReservationRoom
		var playerIDsStr string
		err := rows.Scan(
			&r.ID,
			&r.RoomID,
			&r.Date,
			&playerIDsStr,
		)
		if err != nil {
			return nil, err
		}

		playerIDs := aizuArray(playerIDsStr)

		r.Player, err = searchPlayerInRoom(db, playerIDs)
		if err != nil {
			return nil, err
		}

		reservations = append(reservations, r)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return reservations, nil
}

// insertReservation function
func InsertReservation(db *sql.DB, roomID int, date time.Time) (*int64, error) {
	result, err := db.Exec("INSERT INTO reservations (room_id, date) VALUES (?, ?)", roomID, date)
	if err != nil {
		return nil, err
	}
	id, _ := result.LastInsertId()
	return &id, err

}
