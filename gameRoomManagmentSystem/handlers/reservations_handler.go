package handlers

import (
	"database/sql"
	"net/http"
	"strconv"
	"time"

	"github.com/RYANCOAL9999/SpinnrTechnologyInterview/gameRoomManagmentSystem/databases"
	"github.com/RYANCOAL9999/SpinnrTechnologyInterview/gameRoomManagmentSystem/models"
	"github.com/gin-gonic/gin"
)

const time_format string = "2006-01-02 00:00:00"

func GetReservations(c *gin.Context, db *sql.DB) {
	var args interface{}
	room_id, _ := strconv.Atoi(c.Query("room_id"))
	start_Date, _ := time.Parse(time_format, c.Query("start_Date"))
	end_Date, _ := time.Parse(time_format, c.Query("end_Date"))
	limit, _ := strconv.Atoi(c.Query("limit"))
	reservations, err := databases.ListReservation(db, &room_id, &start_Date, &end_Date, &limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if len(reservations) > 1 {
		args = reservations[0]
	} else {
		args = reservations
	}

	c.JSON(http.StatusOK, args)
}

func updateReservationRoom(db *sql.DB, room_id int, player_ids string) error {
	err := databases.UpdateRoomData(db, room_id, nil, nil, nil, &player_ids)
	if err != nil {
		return err
	}
	return nil
}

func CreateReservations(c *gin.Context, db *sql.DB) {
	var reservation models.ReservationCreate
	if err := c.ShouldBindJSON(&reservation); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	//check the old is available or not
	room, err := databases.ShowRoom(db, reservation.RoomID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if room.Status != models.StatusAvailable {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "this room is not available"})
		return
	}

	id, err := databases.InsertReservation(db, &reservation.RoomID, &reservation.Date)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	var player_ids string = c.Param("player_ids")

	err = updateReservationRoom(db, reservation.ID, player_ids)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"id": id})
}
