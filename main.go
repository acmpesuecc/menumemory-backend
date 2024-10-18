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
	"strconv"
	"time"
)

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
		searchTerm := c.Query("search_term")
		if searchTerm == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "search_term is required"})
			return
		}

		restaurants, err := q.GetRestaurantsLike(c, "%"+searchTerm+"%")
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"restaurants": restaurants,
		})
	})

	r.PUT("/visits/:visit_id", func(c *gin.Context) {
		visitID, err := strconv.ParseInt(c.Param("visit_id"), 10, 64)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid visit_id"})
			return
		}

		userIDStr := c.Query("user_id")
		userID, err := strconv.ParseInt(userIDStr, 10, 64)
		if err != nil || userIDStr == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "user_id is required and should be an integer"})
			return
		}

		var updateData struct {
			Date         string `json:"date"`
			Time         string `json:"time"`
			RestaurantID int64  `json:"restaurant_id"`
		}

		if err := c.ShouldBindJSON(&updateData); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON payload"})
			return
		}

		// Validate if the visit belongs to the user
		visit, err := q.GetVisitByID(c, visitID)
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "visit not found"})
			return
		}

		if visit.Userid.Int64 != userID {
			c.JSON(http.StatusForbidden, gin.H{"error": "You are not authorized to update this visit"})
			return
		}

		// Parse date and time
		parsedDate, err := time.Parse("2006-01-02", updateData.Date)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid date format. Expected YYYY-MM-DD"})
			return
		}

		parsedTime, err := time.Parse("15:04:05", updateData.Time)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid time format. Expected HH:MM:SS"})
			return
		}

		// Update the visit
		err = q.UpdateVisit(c, db.UpdateVisitParams{
			ID:           visitID,
			Date:         parsedDate,
			Time:         parsedTime,
			Restaurantid: sql.NullInt64{Int64: updateData.RestaurantID, Valid: true},
		})
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "Visit updated successfully"})
	})

	return r
}
