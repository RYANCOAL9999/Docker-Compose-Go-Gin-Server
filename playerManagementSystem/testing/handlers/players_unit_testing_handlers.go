package handlers

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	object "github.com/RYANCOAL9999/SpinnrTechnologyInterview/playerManagementSystem/handlers"
	object_models "github.com/RYANCOAL9999/SpinnrTechnologyInterview/playerManagementSystem/models"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestGetPlayers(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Error creating mock database: %v", err)
	}
	defer db.Close()

	rows := sqlmock.NewRows([]string{"ID", "Name", "LV"}).
		AddRow(1, "Alice", 5).
		AddRow(2, "Bob", 10)

	mock.ExpectQuery("SELECT P.ID as ID, P.Name as Name, L.LV as LV FROM Player P").
		WillReturnRows(rows)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	object.GetPlayers(c, db)

	assert.Equal(t, http.StatusOK, w.Code)

	var response []object_models.PlayerRank
	err = json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Len(t, response, 2)
	assert.Equal(t, object_models.PlayerRank{ID: 1, Name: "Alice", LV: 5}, response[0])
	assert.Equal(t, object_models.PlayerRank{ID: 2, Name: "Bob", LV: 10}, response[1])
}

func TestGetPlayers_Error(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Error creating mock database: %v", err)
	}
	defer db.Close()

	mock.ExpectQuery("SELECT P.ID as ID, P.Name as Name, L.LV as LV FROM Player P").
		WillReturnError(sql.ErrConnDone)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	object.GetPlayers(c, db)

	assert.Equal(t, http.StatusInternalServerError, w.Code)

	var response object_models.ErrorResponse
	err = json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Contains(t, response.Error, "error querying database")
}

func TestCreatePlayer(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Error creating mock database: %v", err)
	}
	defer db.Close()

	mock.ExpectExec("INSERT INTO Player").
		WithArgs("Charlie", 3).
		WillReturnResult(sqlmock.NewResult(3, 1))

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	newPlayer := object_models.PlayerRank{Name: "Charlie", LV: 3}
	body, _ := json.Marshal(newPlayer)
	c.Request, _ = http.NewRequest("POST", "/players", bytes.NewBuffer(body))
	c.Request.Header.Set("Content-Type", "application/json")

	object.CreatePlayer(c, db)

	assert.Equal(t, http.StatusCreated, w.Code)

	var response object_models.CreateResponse
	err = json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, 3, response.ID)
}

func TestCreatePlayer_InvalidInput(t *testing.T) {
	db, _, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Error creating mock database: %v", err)
	}
	defer db.Close()

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	invalidPlayer := struct {
		Name string `json:"name"`
		LV   string `json:"lv"` // Invalid type, should be int
	}{
		Name: "Charlie",
		LV:   "not_a_number",
	}
	body, _ := json.Marshal(invalidPlayer)
	c.Request, _ = http.NewRequest("POST", "/players", bytes.NewBuffer(body))
	c.Request.Header.Set("Content-Type", "application/json")

	object.CreatePlayer(c, db)

	assert.Equal(t, http.StatusBadRequest, w.Code)

	var response object_models.ErrorResponse
	err = json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Contains(t, response.Error, "invalid")
}

func TestGetPlayer(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Error creating mock database: %v", err)
	}
	defer db.Close()

	rows := sqlmock.NewRows([]string{"ID", "Name", "LV"}).
		AddRow(1, "Alice", 5)

	mock.ExpectQuery("SELECT P.ID as ID, P.Name as Name, L.LV as LV FROM Player P").
		WithArgs(1).
		WillReturnRows(rows)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Params = []gin.Param{{Key: "id", Value: "1"}}

	object.GetPlayer(c, db)

	assert.Equal(t, http.StatusOK, w.Code)

	var response object_models.PlayerRank
	err = json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, object_models.PlayerRank{ID: 1, Name: "Alice", LV: 5}, response)
}

func TestGetPlayer_Error(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.Default()
	db, _, _ := sqlmock.New()
	defer db.Close()

	router.GET("/player/:id", func(c *gin.Context) {
		object.GetPlayer(c, db)
	})

	req, _ := http.NewRequest(http.MethodGet, "/player/abc", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
}

func TestUpdatePlayer(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Error creating mock database: %v", err)
	}
	defer db.Close()

	mock.ExpectQuery("SELECT ID FROM Level WHERE LV = ?").
		WithArgs(6).
		WillReturnRows(sqlmock.NewRows([]string{"ID"}).AddRow(2))

	mock.ExpectExec("UPDATE players SET").
		WithArgs(2, "Alice Updated", 1).
		WillReturnResult(sqlmock.NewResult(0, 1))

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	updatedPlayer := object_models.PlayerRank{ID: 1, Name: "Alice Updated", LV: 6}
	body, _ := json.Marshal(updatedPlayer)
	c.Request, _ = http.NewRequest("PUT", "/players", bytes.NewBuffer(body))
	c.Request.Header.Set("Content-Type", "application/json")

	object.UpdatePlayer(c, db)

	assert.Equal(t, http.StatusOK, w.Code)

	var response object_models.SuccessResponse
	err = json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
}

func TestUpdatePlayer_Error(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.Default()
	db := &sql.DB{}

	router.PUT("/player", func(c *gin.Context) {
		object.UpdatePlayer(c, db)
	})

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodPut, "/player", bytes.NewBuffer([]byte("invalid json")))
	req.Header.Set("Content-Type", "application/json")

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestDeletePlayer(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Error creating mock database: %v", err)
	}
	defer db.Close()

	mock.ExpectExec("DELETE FROM players WHERE ID = ?").
		WithArgs(1).
		WillReturnResult(sqlmock.NewResult(0, 1))

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Params = []gin.Param{{Key: "id", Value: "1"}}

	object.DeletePlayer(c, db)

	assert.Equal(t, http.StatusOK, w.Code)

	var response object_models.SuccessResponse
	err = json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
}

func TestDeletePlayer_Error(t *testing.T) {
	gin.SetMode(gin.TestMode)
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	mock.ExpectExec("DELETE FROM players WHERE id = ?").
		WithArgs(999).
		WillReturnError(sql.ErrNoRows)

	req, _ := http.NewRequest(http.MethodDelete, "/players/999", nil)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = req
	c.Params = gin.Params{gin.Param{Key: "id", Value: "999"}}

	object.DeletePlayer(c, db)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
	assert.NoError(t, mock.ExpectationsWereMet())
}
