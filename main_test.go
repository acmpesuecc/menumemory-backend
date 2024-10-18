package main

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

var router *gin.Engine

func TestMain(m *testing.M) {
	router = SetupApp()
	m.Run()
}

func TestPingRoute(t *testing.T) {
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/ping", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, 200, w.Code)
	assert.Equal(t, "Pong Uwu", w.Body.String())
}

func TestGetRestaurants(t *testing.T) {
	// Test not passing in "search_term" to get 400 error
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/restaurants", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, 400, w.Code)

	// Test passing "search_term" = "Milano"
	w = httptest.NewRecorder()
	req, _ = http.NewRequest("GET", "/restaurants?search_term=Milano", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, 200, w.Code)

	var response struct {
		Restaurants []struct {
			Name string `json:"Name"`
		} `json:"restaurants"`
	}

	err := json.NewDecoder(w.Body).Decode(&response)
	assert.NoError(t, err)

	// Assert that each restaurant's name contains "Milano"
	for _, restaurant := range response.Restaurants {
		assert.Contains(t, strings.ToLower(restaurant.Name), "milano")
	}
}

func TestUpdateVisit(t *testing.T) {
	// Mock data for the update
	visitID := "1" // Example visit ID
	userID := "1"   // Example user ID (must match the user in the DB for the test)

	// Create request body
	requestBody := `{
		"date": "2021-10-17",
		"time": "18:24:00",
		"restaurant_id": 1
	}`

	// Create a recorder and a request
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("PUT", "/visits/"+visitID+"?user_id="+userID, strings.NewReader(requestBody))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	// Check for success response
	assert.Equal(t, 200, w.Code)
	assert.Contains(t, w.Body.String(), "Visit updated successfully")

	// Test unauthorized access
	w = httptest.NewRecorder()
	req, _ = http.NewRequest("PUT", "/visits/"+visitID+"?user_id=2", strings.NewReader(requestBody))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	// Check for forbidden response
	assert.Equal(t, 403, w.Code)
	assert.Contains(t, w.Body.String(), "Unauthorized access")

	// Test invalid request body
	w = httptest.NewRecorder()
	req, _ = http.NewRequest("PUT", "/visits/"+visitID+"?user_id="+userID, strings.NewReader(`{"invalid": "data"}`))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	// Check for bad request response
	assert.Equal(t, 400, w.Code)
	assert.Contains(t, w.Body.String(), "Invalid request body")
}
