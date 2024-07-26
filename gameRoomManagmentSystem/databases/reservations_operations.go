package databases

import (
	"database/sql"
	"fmt"
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
		return nil, fmt.Errorf("error querying database with ListReservation: %w", err)
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
			return nil, fmt.Errorf("error scanning row with ListReservation: %w", err)
		}

		playerIDs := aizuArray(playerIDsStr)

		r.Player, _ = searchPlayerInRoom(db, playerIDs)

		reservations = append(reservations, r)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating over rows with ListReservation: %w", err)
	}

	return reservations, nil
}

// insertReservation function
func InsertReservation(db *sql.DB, roomID int, date time.Time) (int, error) {
	result, err := db.Exec("INSERT INTO reservations (room_id, date) VALUES (?, ?)", roomID, date)
	if err != nil {
		return 0, fmt.Errorf("error querying database with InsertReservation: %w", err)
	}
	id, _ := result.LastInsertId()
	return int(id), err

}
