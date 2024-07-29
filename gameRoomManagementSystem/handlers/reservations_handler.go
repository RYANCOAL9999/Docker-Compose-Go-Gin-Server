package handlers

import (
	"database/sql"
	"net/http"
	"strconv"
	"time"

	"github.com/RYANCOAL9999/SpinnrTechnologyInterview/gameRoomManagementSystem/databases"
	"github.com/RYANCOAL9999/SpinnrTechnologyInterview/gameRoomManagementSystem/models"
	"github.com/gin-gonic/gin"
)

const time_format string = "2006-01-02 00:00:00"

// @Summary      Retrieve reservations
// @Description  Get a list of reservations based on optional filters such as room ID, start date, end date, and limit. Returns reservations that match the criteria.
// @Tags         reservations
// @Accept       json
// @Produce      json
// @Param        room_id    query  int     false  "Filter reservations by room ID"
// @Param        start_date query  string  false  "Filter reservations starting from this date (format: YYYY-MM-DD)"
// @Param        end_date   query  string  false  "Filter reservations ending at this date (format: YYYY-MM-DD)"
// @Param        limit      query  int     false  "Limit the number of results returned"
// @Success      200  {object}  []models.ReservationRoom  "List of reservations matching the criteria"
// @Failure      400  {object}  models.ErrorResponse "Bad request due to invalid query parameters"
// @Failure      500  {object}  models.ErrorResponse "Internal server error"
// @Router       /reservations [get]
func GetReservations(c *gin.Context, db *sql.DB) {
	var args interface{}
	roomID, _ := strconv.Atoi(c.Query("room_id"))
	startDate, _ := time.Parse(time_format, c.Query("start_Date"))
	endDate, _ := time.Parse(time_format, c.Query("end_Date"))
	limit, _ := strconv.Atoi(c.Query("limit"))
	reservations, err := databases.ListReservation(db, roomID, startDate, endDate, limit)
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

func UpdateReservationRoom(db *sql.DB, roomID int, playerIDs string) error {
	var room models.Room
	room.ID = roomID
	room.PlayerIDs = playerIDs
	err := databases.UpdateRoomData(db, room)
	if err != nil {
		return err
	}
	return nil
}

// CreateReservations handles the creation of a new reservation for a room.
// It processes the request to create a reservation by checking if the room is available,
// validating the date format, inserting the reservation into the database, and updating the room with player IDs.
//
// @Summary      Create a reservation
// @Description  Creates a new reservation for a specified room if the room is available. The request body must include the room ID, date of reservation, and player IDs. If successful, returns the ID of the created reservation.
// @Tags         reservations
// @Accept       json
// @Produce      json
// @Param        reservation  body  models.Reservation  true  "Reservation details to be created"
// @Success      201  {object}  models.CreateResponse "Reservation created successfully, returns the ID of the new reservation"
// @Failure      400  {object}  models.ErrorResponse "Bad request due to invalid input or date format"
// @Failure      500  {object}  models.ErrorResponse "Internal server error"
// @Router       /reservations [post]
func CreateReservations(c *gin.Context, db *sql.DB) {
	var reservation models.Reservation
	if err := c.ShouldBindJSON(&reservation); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	//check the old is available or not
	room, err := databases.ShowRoom(db, reservation.RoomID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{Error: err.Error()})
		return
	}

	if room.Status != models.StatusAvailable {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{Error: "this room is not available"})
		return
	}

	time, err := time.Parse(time_format, reservation.Date)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{Error: "date format has issues"})
		return
	}

	id, err := databases.InsertReservation(db, reservation.RoomID, time)

	if err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{Error: err.Error()})
		return
	}

	err = UpdateReservationRoom(db, reservation.RoomID, reservation.PlayerIDs)

	if err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusCreated, models.CreateResponse{ID: id})
}
