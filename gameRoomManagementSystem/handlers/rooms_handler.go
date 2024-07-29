package handlers

import (
	"database/sql"
	"net/http"
	"strconv"

	"github.com/RYANCOAL9999/SpinnrTechnologyInterview/gameRoomManagementSystem/databases"
	"github.com/RYANCOAL9999/SpinnrTechnologyInterview/gameRoomManagementSystem/models"

	"github.com/gin-gonic/gin"
)

// @Summary      Retrieve all rooms
// @Description  Get a list of all rooms available in the database along with their details such as name, status, description, and player IDs.
// @Tags         rooms
// @Accept       json
// @Produce      json
// @Success      200  {object}  []models.Room  "List of rooms"
// @Failure      500  {object}  models.ErrorResponse  "Internal server error"
// @Router       /rooms [get]
func GetRooms(c *gin.Context, db *sql.DB) {
	rooms, err := databases.ListRooms(db)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{Error: err.Error()})
		return
	}
	c.JSON(http.StatusOK, rooms)
}

// @Summary      Create a new room
// @Description  Add a new room to the database with the provided name, description, and status. PlayerIDs are optional and can be set later.
// @Tags         rooms
// @Accept       json
// @Produce      json
// @Param        room  body  models.Room  true  "Room details to be created"
// @Success      201  {object}  models.CreateResponse "ID of the created room"
// @Failure      400  {object}  models.ErrorResponse  "Bad request due to invalid input"
// @Failure      500  {object}  models.ErrorResponse  "Internal server error"
// @Router       /rooms [post]
func CreateRoom(c *gin.Context, db *sql.DB) {
	var room models.Room
	if err := c.BindJSON(&room); err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{Error: err.Error()})
		return
	}

	id, err := databases.AddRoom(db, room.Name, room.Description)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{Error: err.Error()})
		return
	}
	c.JSON(http.StatusCreated, models.CreateResponse{ID: id})
}

// @Summary      Retrieve a room by ID
// @Description  Fetch details of a specific room from the database identified by its ID. Returns room details including name, status, description, and player IDs.
// @Tags         rooms
// @Accept       json
// @Produce      json
// @Param        id  path  int  true  "Room ID"
// @Success      200  {object}  []models.Room  "Details of the room"
// @Failure      400  {object}  models.ErrorResponse "Invalid ID supplied"
// @Failure      500  {object}  models.ErrorResponse "Internal server error"
// @Router       /rooms/{id} [get]
func GetRoom(c *gin.Context, db *sql.DB) {
	id, _ := strconv.Atoi(c.Param("id"))
	room, err := databases.ShowRoom(db, id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{Error: err.Error()})
		return
	}
	c.JSON(http.StatusOK, room)
}

// @Summary      Update a room
// @Description  Update the details of an existing room in the database. The request body should include the room's ID, name, status, description, and player IDs. The ID is used to identify the room to be updated.
// @Tags         rooms
// @Accept       json
// @Produce      json
// @Param        room  body  models.Room  true  "Room details to be updated"
// @Success      200  {object}  models.SuccessResponse "Update successful"
// @Failure      400  {object}  models.ErrorResponse   "Bad request due to invalid input"
// @Failure      500  {object}  models.ErrorResponse   "Internal server error"
// @Router       /rooms [put]
func UpdateRoom(c *gin.Context, db *sql.DB) {
	var room models.Room
	if err := c.BindJSON(&room); err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{Error: err.Error()})
		return
	}
	err := databases.UpdateRoomData(db, room)

	if err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{Error: err.Error()})
		return
	}
	c.JSON(http.StatusOK, models.SuccessResponse{})
}

// @Summary      Delete a room
// @Description  Remove a specific room from the database using its ID. If the room exists, it will be deleted.
// @Tags         rooms
// @Accept       json
// @Produce      json
// @Param        id  path  int  true  "Room ID"
// @Success      200  {object}  models.SuccessResponse "Delete successful"
// @Failure      400  {object}  models.ErrorResponse   "Invalid ID supplied"
// @Failure      500  {object}  models.ErrorResponse   "Internal server error"
// @Router       /rooms/{id} [delete]
func DeleteRoom(c *gin.Context, db *sql.DB) {
	id, _ := strconv.Atoi(c.Param("id"))
	err := databases.DeleteRoom(db, id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{Error: err.Error()})
		return
	}
	c.JSON(http.StatusOK, models.SuccessResponse{})
}
