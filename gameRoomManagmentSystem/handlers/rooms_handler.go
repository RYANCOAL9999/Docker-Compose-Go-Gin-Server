package handlers

import (
	"database/sql"
	"net/http"
	"strconv"

	"github.com/RYANCOAL9999/SpinnrTechnologyInterview/gameRoomManagmentSystem/databases"
	"github.com/RYANCOAL9999/SpinnrTechnologyInterview/gameRoomManagmentSystem/models"

	"github.com/gin-gonic/gin"
)

func GetRooms(c *gin.Context, db *sql.DB) {
	rooms, err := databases.ListRooms(db)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, rooms)
}

func CreateRoom(c *gin.Context, db *sql.DB) {
	var room models.Room
	if err := c.BindJSON(&room); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	id, err := databases.AddRoom(db, room.Name, room.Description)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"id": id})
}

func GetRoom(c *gin.Context, db *sql.DB) {
	id, _ := strconv.Atoi(c.Param("id"))
	room, err := databases.ShowRoom(db, id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, room)
}

func UpdateRoom(c *gin.Context, db *sql.DB) {
	var room models.Room
	if err := c.BindJSON(&room); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := databases.UpdateRoomData(db, room.ID, &room.Name, (*int)(&room.Status), &room.Description, &room.Player_ids)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err})
		return
	}
	var args interface{}
	c.JSON(http.StatusOK, args)
}

func DeleteRoom(c *gin.Context, db *sql.DB) {
	id, _ := strconv.Atoi(c.Param("id"))
	err := databases.DeleteRoom(db, id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err})
		return
	}

	var args interface{}
	c.JSON(http.StatusOK, args)
}
