package databases

import (
	"database/sql"
	"fmt"
	"strings"

	"github.com/RYANCOAL9999/SpinnrTechnologyInterview/gameRoomManagmentSystem/models"
)

func ListRooms(db *sql.DB) ([]models.Room, error) {
	var query string = `
		SELECT 
		id, name, status 
		FROM rooms
	`

	rows, err := db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("error querying database with ListRooms: %w", err)
	}
	defer rows.Close()

	var rooms []models.Room
	for rows.Next() {
		var room models.Room
		err := rows.Scan(
			&room.ID,
			&room.Name,
			&room.Status,
		)
		if err != nil {
			return nil, fmt.Errorf("error scanning row with ListRooms: %w", err)
		}
		rooms = append(rooms, room)
	}
	return rooms, nil
}

func ShowRoom(db *sql.DB, id int) (*models.Room, error) {
	var room models.Room
	err := db.QueryRow(`
		SELECT 
		id, name, status 
		FROM players 
		WHERE id = ?
	`, id).Scan(
		&room.ID,
		&room.Name,
		&room.Status,
	)
	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("error querying database with Show Room: %w", err)
	} else if err != nil {
		return nil, fmt.Errorf("error scanning row with Show Room: %w", err)
	}
	return &room, nil
}

func AddRoom(db *sql.DB, name string, description string) (int, error) {
	result, err := db.Exec(`
		INSERT INTO room (name, status, description) 
		VALUES (?, ?, ?)
	`, name, models.StatusAvailable, description)
	if err != nil {
		return 0, fmt.Errorf("error querying database with AddRoom: %w", err)
	}
	id, _ := result.LastInsertId()
	return int(id), err
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
		return fmt.Errorf("no fields update for Room with id: %d", room.ID)
	}

	query += " " + strings.Join(updates, ", ")
	query += " WHERE id = ?"
	args = append(args, room.ID)

	result, err := db.Exec(query, args...)
	if err != nil {
		return fmt.Errorf("error querying database with UpdateRoomData: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("error updating room: %w", err)
	}
	if rowsAffected == 0 {
		return fmt.Errorf("no rows were updated, room with id %d may not exist", room.ID)
	}

	return nil
}

func DeleteRoom(db *sql.DB, id int) error {
	result, err := db.Exec(`
		DELETE FROM players 
		WHERE id = ?
	`, id)
	if err != nil {
		return fmt.Errorf("error querying database with DeleteRoom: %w", err)
	}
	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		return fmt.Errorf("room not found")
	}

	return nil
}

func searchPlayerInRoom(db *sql.DB, playerIDs []int) ([]models.PlayerRank, error) {
	// Create a string of placeholders for the IN clause
	placeholders := make([]string, len(playerIDs))
	args := make([]interface{}, len(playerIDs))
	for i, id := range playerIDs {
		placeholders[i] = "?"
		args[i] = id
	}

	var query string = fmt.Sprintf(`
        SELECT 
		p.id, p.name, l.rank 
        FROM players p
        INNER JOIN levels l ON p.level_id = l.id
        WHERE p.id IN (%s)
        ORDER BY p.id
    `, strings.Join(placeholders, ","))

	rows, err := db.Query(query, args...)
	if err != nil {
		return nil, fmt.Errorf("error querying database with searchPlayerInRoom: %w", err)
	}
	defer rows.Close()

	var playerRanks []models.PlayerRank
	for rows.Next() {
		var playerRank models.PlayerRank
		err := rows.Scan(
			&playerRank.ID,
			&playerRank.Name,
			&playerRank.Rank,
		)
		if err != nil {
			return nil, fmt.Errorf("error scanning row with searchPlayerInRoom: %w", err)
		}
		playerRanks = append(playerRanks, playerRank)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating over rows with searchPlayerInRoom: %w", err)
	}

	// Check if all requested players were found
	if len(playerRanks) != len(playerIDs) {
		return playerRanks, fmt.Errorf("some players were not found")
	}

	return playerRanks, nil
}
