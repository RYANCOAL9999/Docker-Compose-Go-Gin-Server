package databases

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/RYANCOAL9999/SpinnrTechnologyInterview/gameLogCollector/models"
)

func ListLogs(db *sql.DB, playerID int, action string, startTime time.Time, endTime time.Time, limit int) ([]models.Log, error) {

	query := `SELECT id, player_id, action, timestamp, details FROM game_logs WHERE 1=1`
	var args []interface{}

	if playerID != 0 {
		query += " AND player_id = ?"
		args = append(args, playerID)
	}
	if action != "" {
		query += " AND action = ?"
		args = append(args, action)
	}

	if !startTime.IsZero() && !endTime.IsZero() {
		query += " AND timestamp BETWEEN ? AND ?"
		args = append(args, startTime, endTime)
	}

	query += " ORDER BY timestamp DESC"

	if limit != 0 {
		query += " LIMIT ?"
		args = append(args, limit)
	}

	rows, err := db.Query(query, args...)
	if err != nil {
		return nil, fmt.Errorf("error querying database with ListLogs: %w", err)
	}
	defer rows.Close()

	var logs []models.Log
	for rows.Next() {
		var log models.Log
		err := rows.Scan(&log.ID, &log.PlayerID, &log.Action, &log.Timestamp, &log.Details)
		if err != nil {
			return nil, fmt.Errorf("error scanning row with with ListLogs: %w", err)
		}
		logs = append(logs, log)
	}
	return logs, nil
}

func AddLog(db *sql.DB, log models.Log) (int, error) {
	result, err := db.Exec(`
		INSERT INTO game_logs (player_id, action, timestamp, details) VALUES (?, ?, ?, ?)
	`, log.PlayerID, log.Action, time.Now(), log.Details)
	if err != nil {
		return 0, fmt.Errorf("error querying database with AddLog: %w", err)
	}
	id, _ := result.LastInsertId()
	return int(id), err
}
