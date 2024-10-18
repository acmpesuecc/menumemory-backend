package main

import (
	"bytes"
	"context"
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
	"github.com/stretchr/testify/assert"
	"your_project/db" // Replace with the actual path to the db package
)

var testQueries *db.Queries
var testDB *sql.DB

// Setup the database for testing
func setup() {
	var err error
	testDB, err = sql.Open("postgres", "postgresql://user:password@localhost:5432/your_test_db?sslmode=disable")
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}
	testQueries = db.New(testDB)
}

// Teardown the database connection after testing
func teardown() {
	testDB.Close()
}

// TestMain is used to set up and tear down the testing environment.
func TestMain(m *testing.M) {
	gin.SetMode(gin.TestMode)
	setup()
	defer teardown()
	m.Run()
}

// Helper function to create a new HTTP request and return the response recorder.
func executeRequest(req *http.Request) *httptest.ResponseRecorder {
	rr := httptest.NewRecorder()
	router := setupRouter() // Assuming you have a function to setup your router.
	router.ServeHTTP(rr, req)
	return rr
}

// TestCreateVisit tests the creation of a new visit.
func TestCreateVisit(t *testing.T) {
	params := db.CreateVisitParams{
		Date:         time.Now(),
		Time:         "12:30 PM",
		Userid:       sql.NullInt64{Int64: 1, Valid: true},
		Restaurantid: sql.NullInt64{Int64: 1, Valid: true},
	}

	err := testQueries.CreateVisit(context.Background(), params)
	assert.NoError(t, err)
}

// TestGetVisitById tests fetching a visit by its ID.
func TestGetVisitById(t *testing.T) {
	visitID := int64(1)
	visit, err := testQueries.GetVisitById(context.Background(), visitID)

	assert.NoError(t, err)
	assert.Equal(t, visit.ID, visitID)
	assert.NotEmpty(t, visit.Date)
	assert.NotEmpty(t, visit.Time)
}

// TestGetOrdersForVisit tests fetching orders for a specific visit.
func TestGetOrdersForVisit(t *testing.T) {
	visitID := sql.NullInt64{Int64: 1, Valid: true}
	orders, err := testQueries.GetOrdersForVisit(context.Background(), visitID)

	assert.NoError(t, err)
	assert.NotNil(t, orders)
}

// TestCreateOrder tests creating a new order.
func TestCreateOrder(t *testing.T) {
	params := db.CreateOrderParams{
		Visitid:    sql.NullInt64{Int64: 1, Valid: true},
		Dishid:     sql.NullInt64{Int64: 1, Valid: true},
		Rating:     sql.NullFloat64{Float64: 4.5, Valid: true},
		Reviewtext: sql.NullString{String: "Delicious dish!", Valid: true},
	}

	err := testQueries.CreateOrder(context.Background(), params)
	assert.NoError(t, err)
}

// TestGetRestaurantsLike tests fetching restaurants with a name pattern.
func TestGetRestaurantsLike(t *testing.T) {
	name := "%Pizza%"
	restaurants, err := testQueries.GetRestaurantsLike(context.Background(), name)

	assert.NoError(t, err)
	assert.NotEmpty(t, restaurants)
}

// TestUpdateVisit tests updating an existing visit.
func TestUpdateVisit(t *testing.T) {
	params := db.UpdateVisitParams{
		Date:         time.Now(),
		Time:         "01:00 PM",
		Restaurantid: sql.NullInt64{Int64: 2, Valid: true},
		ID:           1,
		Userid:       sql.NullInt64{Int64: 1, Valid: true},
	}

	err := testQueries.UpdateVisit(context.Background(), params)
	assert.NoError(t, err)
}

// TestAPIGetVisitByID tests the API endpoint for fetching a visit by its ID.
func TestAPIGetVisitByID(t *testing.T) {
	req, _ := http.NewRequest("GET", "/visits/1", nil)
	response := executeRequest(req)

	assert.Equal(t, http.StatusOK, response.Code)

	var visit db.GetVisitByIdRow
	err := json.Unmarshal(response.Body.Bytes(), &visit)
	assert.NoError(t, err)
	assert.Equal(t, int64(1), visit.ID)
}

// TestAPICreateVisit tests the API endpoint for creating a new visit.
func TestAPICreateVisit(t *testing.T) {
	payload := map[string]interface{}{
		"date":          time.Now().Format("2006-01-02"),
		"time":          "02:00 PM",
		"userId":        1,
		"restaurantId":  1,
	}
	jsonPayload, _ := json.Marshal(payload)

	req, _ := http.NewRequest("POST", "/visits", bytes.NewBuffer(jsonPayload))
	req.Header.Set("Content-Type", "application/json")
	response := executeRequest(req)

	assert.Equal(t, http.StatusCreated, response.Code)
}

// TestAPIGetOrdersForVisit tests the API endpoint for fetching orders for a visit.
func TestAPIGetOrdersForVisit(t *testing.T) {
	req, _ := http.NewRequest("GET", "/visits/1/orders", nil)
	response := executeRequest(req)

	assert.Equal(t, http.StatusOK, response.Code)

	var orders []db.GetOrdersForVisitRow
	err := json.Unmarshal(response.Body.Bytes(), &orders)
	assert.NoError(t, err)
	assert.NotEmpty(t, orders)
}
