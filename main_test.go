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
	//Test not passing in "search_term" to get 400 erro
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

	//t.Log(response.Restaurants)

	// Assert that each restaurant's name contains "Milano"
	for _, restaurant := range response.Restaurants {
		assert.Contains(t, strings.ToLower(restaurant.Name), "milano")
	}
}


func TestDeleteVisitHandler(t *testing.T) {
    // Test case 1: Missing user_id
    w := httptest.NewRecorder()
    req, _ := http.NewRequest("DELETE", "/visits/1", nil)
    router.ServeHTTP(w, req)
    assert.Equal(t, 400, w.Code)
    assert.Contains(t, w.Body.String(), "user_id is required")

    // Test case 2: Unauthorized access (user_id != "1")
    w = httptest.NewRecorder()
    req, _ = http.NewRequest("DELETE", "/visits/1?user_id=2", nil)
    router.ServeHTTP(w, req)
    assert.Equal(t, 403, w.Code)
    assert.Contains(t, w.Body.String(), "Unauthorized access")

    // Test case 3: Invalid visit_id
    w = httptest.NewRecorder()
    req, _ = http.NewRequest("DELETE", "/visits/invalid?user_id=1", nil)
    router.ServeHTTP(w, req)
    assert.Equal(t, 400, w.Code)
    assert.Contains(t, w.Body.String(), "Invalid visit_id")

    // Test case 4: Visit not found
    // This assumes that visit with ID 9999 doesn't exist in your test database
    w = httptest.NewRecorder()
    req, _ = http.NewRequest("DELETE", "/visits/9999?user_id=1", nil)
    router.ServeHTTP(w, req)
    assert.Equal(t, 404, w.Code)
    assert.Contains(t, w.Body.String(), "Visit not found")

    // Test case 5: Visit doesn't belong to user
    // This assumes that visit with ID 2 exists but doesn't belong to user 1
    w = httptest.NewRecorder()
    req, _ = http.NewRequest("DELETE", "/visits/2?user_id=1", nil)
    router.ServeHTTP(w, req)
    assert.Equal(t, 403, w.Code)
    assert.Contains(t, w.Body.String(), "Visit does not belong to user")

    // Test case 6: Successful deletion
    // This assumes that visit with ID 1 exists and belongs to user 1
    w = httptest.NewRecorder()
    req, _ = http.NewRequest("DELETE", "/visits/1?user_id=1", nil)
    router.ServeHTTP(w, req)
    assert.Equal(t, 204, w.Code)
    assert.Empty(t, w.Body.String())
}