package databases

import (
	"database/sql"
	"errors"

	"github.com/RYANCOAL9999/SpinnrTechnologyInterview/gameRoomManagmentSystem/models"
)

func ListRooms(db *sql.DB) ([]models.Room, error) {
	var query string = "SELECT id, name, status FROM rooms"

	rows, err := db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var rooms []models.Room
	for rows.Next() {
		var room models.Room
		if err := rows.Scan(&room.ID, &room.Name, &room.Status); err != nil {
			return nil, err
		}
		rooms = append(rooms, room)
	}
	return rooms, nil
}

func ShowRoom(db *sql.DB, id int) (*models.Room, error) {
	var room models.Room
	err := db.QueryRow("SELECT id, name, status FROM players WHERE id = ?", id).Scan(&room.ID, &room.Name, &room.Status)
	if err == sql.ErrNoRows {
		return nil, err
	} else if err != nil {
		return nil, err
	}
	return &room, nil
}

func AddRoom(db *sql.DB, name string, description string) (*int64, error) {
	result, err := db.Exec("INSERT INTO room (name, status, description) VALUES (?, ?)", name, models.StatusAvailable, description)
	if err != nil {
		return nil, err
	}
	id, _ := result.LastInsertId()
	return &id, err
}

func UpdateRoomData(db *sql.DB, id int, name *string, status *int, description *string, player_ids *string) error {
	var _ sql.Result
	var err error

	var query string = `UPDATE rooms SET`
	args := []interface{}{}

	if name != nil {
		query += " name = ?"
		args = append(args, name)
	}

	if status != nil {
		query += " status = ?"
		args = append(args, status)
	}

	if description != nil {
		query += " description = ?"
		args = append(args, description)
	}

	if player_ids != nil {
		query += " player_ids = ?"
		args = append(args, player_ids)
	}

	args = append(args, id)

	_, err = db.Exec(query, args...)

	if err != nil {
		return err
	}
	return nil
}

func DeleteRoom(db *sql.DB, id int) error {
	result, err := db.Exec("DELETE FROM players WHERE id = ?", id)
	if err != nil {
		return err
	}
	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		err = errors.New("room not found")
		return err
	}

	return nil
}

func searchPlayerInRoom(db *sql.DB, player_id []int) ([]models.PlayerRank, error) {
	rows, err := db.Query(`
		SELECT p.id, p.name, l.rank 
		FROM players p
		INNER JOIN 
		level l 
		ON p.level_id = l.id
		WHERE p.id in (?)
		ORDER BY p.id
	`, player_id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var playerRanks []models.PlayerRank
	for rows.Next() {
		var playerRank models.PlayerRank
		if err := rows.Scan(&playerRank.ID, &playerRank.Name, &playerRank.Rank); err != nil {
			return nil, err
		}
		playerRanks = append(playerRanks, playerRank)
	}
	return playerRanks, nil
}
