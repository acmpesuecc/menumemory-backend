package main

import (
	"database/sql"
	"fmt"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	_ "github.com/mattn/go-sqlite3"
	"menumemory-backend/db"
	"net/http"
	"os"
)

// Visit struct to represent visit data
type Visit struct {
	ID          string // Add the Visit ID field
	UserID      string // Assuming this field exists to track ownership
	Date        string // Visit date
	Time        string // Visit time
	RestaurantID int    // Associated restaurant ID
}

func SetupApp() *gin.Engine {
	fmt.Println("Beginning Database Initialization")
	db_, err := sql.Open("sqlite3", "warehouse.db")
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
	q := db.New(db_)
	fmt.Println("Finished Database Initialization")

	r := gin.Default()

	// Allow all cors origins
	r.Use(cors.Default())

	r.StaticFile("/openapi.json", "./openapi.json")
	r.StaticFile("/openapi.yaml", "./openapi.yaml")

	r.GET("/ping", func(c *gin.Context) {
		c.String(200, "Pong Uwu")
	})

	r.GET("/restaurants", func(c *gin.Context) {
		search_term := c.Query("search_term")
		if search_term == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "search_term is required"})
		}

		restaurants, err := q.GetRestaurantsLike(c, "%"+search_term+"%")
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}

		c.JSON(http.StatusOK, gin.H{
			"restaurants": restaurants,
		})
	})

	// Add the PUT endpoint for updating visits
	r.PUT("/visits/:visit_id", func(c *gin.Context) {
		visitID := c.Param("visit_id")
		userID := c.Query("user_id")

		if userID == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "user_id is required"})
			return
		}

		var req struct {
			Date         string `json:"date" binding:"required"`
			Time         string `json:"time" binding:"required"`
			RestaurantID int    `json:"restaurant_id" binding:"required"`
		}

		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
			return
		}

		// Check if the visit belongs to the user
		visit, err := getVisitByID(visitID) // Implement this function to fetch visit details
		if err != nil || visit.UserID != userID {
			c.JSON(http.StatusForbidden, gin.H{"error": "Unauthorized"})
			return
		}

		// Update the visit in the database
		err = updateVisitInDB(visitID, req.Date, req.Time, req.RestaurantID) // Implement this function
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to update visit"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "visit updated successfully"})
	})

	return r
}

func main() {
	r := SetupApp()
	r.Run() // listen and serve on 0.0.0.0:8080
}

// Add functions to interact with the database
func getVisitByID(visitID string) (*Visit, error) {
	// Implement your logic to retrieve the visit by ID from the database
	return &Visit{}, nil // Placeholder return
}

func updateVisitInDB(visitID, date, time string, restaurantID int) error {
	// Implement your logic to update the visit in the database
	return nil // Placeholder return
}
