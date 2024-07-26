package databases

import (
	"database/sql"
	"errors"
	"strings"

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

func UpdateRoomData(db *sql.DB, room models.Room) error {
	query := "UPDATE rooms SET"
	args := []interface{}{}
	updates := []string{}

	if room.Name != "" {
		updates = append(updates, "name = ?")
		args = append(args, room.Name)
	}

	if room.Status != 0 {
		updates = append(updates, "status = ?")
		args = append(args, room.Status)
	}

	if room.Description != "" {
		updates = append(updates, "description = ?")
		args = append(args, room.Description)
	}

	if room.PlayerIDs != "" {
		updates = append(updates, "player_ids = ?")
		args = append(args, room.PlayerIDs)
	}

	if len(updates) == 0 {
		return errors.New("no fields to update")
	}

	query += " " + strings.Join(updates, ", ")
	query += " WHERE id = ?"
	args = append(args, room.ID)

	result, err := db.Exec(query, args...)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return errors.New("no rows were updated")
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
		return errors.New("room not found")
	}

	return nil
}

func searchPlayerInRoom(db *sql.DB, playerIDs []int) ([]models.PlayerRank, error) {
	rows, err := db.Query(`
		SELECT p.id, p.name, l.rank 
		FROM players p
		INNER JOIN 
		level l 
		ON p.level_id = l.id
		WHERE p.id in (?)
		ORDER BY p.id
	`, playerIDs)
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
