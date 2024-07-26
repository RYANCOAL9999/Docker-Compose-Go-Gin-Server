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

func ListReservation(db *sql.DB, room_id *int, start_Date *time.Time, end_Date *time.Time, limit *int) ([]models.Reservation, error) {
	var query string = `
		SELECT
		r.id AS reservation_id,
		rm.id AS room_id,
		r.date_time AS reservation_date_time,
		rm.Player_ids AS player_ids,
		FROM reservations r
		INNER JOIN 
		room rm 
		ON r.room_id = rm.id
	`
	args := []interface{}{}

	if room_id != nil {
		query += " AND room_id = ?"
		args = append(args, room_id)
	}

	if start_Date != nil && end_Date != nil {
		query += " AND date Between ? "
		args = append(args, start_Date)
		query += " AND ? "
		args = append(args, end_Date)
	}

	query += `ORDER BY reservation_id`

	if limit != nil {
		query += " LIMIT ?"
		args = append(args, limit)
	}

	rows, err := db.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var reservation []models.Reservation
	for rows.Next() {
		var r models.Reservation
		var longStr string
		err := rows.Scan(
			&r.ID,
			&r.RoomID,
			&r.Date,
			&longStr,
		)
		if err != nil {
			return nil, err
		}

		player_Ids := aizuArray(longStr)

		r.Player, err = searchPlayerInRoom(db, player_Ids)

		if err != nil {
			return nil, err
		}

		reservation = append(reservation, r)
	}

	return reservation, err
}

// insertReservation function
func InsertReservation(db *sql.DB, room_id *int, date *time.Time) (*int64, error) {
	//use the room_id to check need to check the is not resveration
	result, err := db.Exec("INSERT INTO reservations (room_id, date, player_id) VALUES (?, ?, ?)", room_id, date)
	if err != nil {
		return nil, err
	}
	id, _ := result.LastInsertId()
	return &id, err

}
