package main

import (
	"database/sql"
	"fmt"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"menumemory-backend/db"
	"net/http"
	_ "github.com/mattn/go-sqlite3"
	"os"
	"strconv"
)

var q *db.Queries

func SetupApp() *gin.Engine {
	fmt.Println("Beginning Database Initialization")
	db_, err := sql.Open("sqlite3", "warehouse.db")
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
	// Assign the global variable
	q = db.New(db_)
	fmt.Println("Finished Database Initialization")

	r := gin.Default()

	// Allow all CORS origins
	r.Use(cors.Default())

	r.DELETE("/visits/:visit_id", DeleteVisitHandler)

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

	return r
}
func DeleteVisitHandler(c *gin.Context) {
    visitID := c.Param("visit_id")
    userID := c.Query("user_id")

    if userID == "" {
        c.JSON(http.StatusBadRequest, gin.H{"error": "user_id is required"})
        return
    }

    // Check if userID is valid
    if userID != "1" {
        c.JSON(http.StatusForbidden, gin.H{"error": "Unauthorized access"})
        return
    }

    // Convert visitID to an integer
    visitIDInt, err := strconv.Atoi(visitID)
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid visit_id"})
        return
    }

    // Check if visit belongs to user_id
    visitOwner, err := q.GetVisitOwnerByID(c, int64(visitIDInt))
    if err != nil {
        if err == sql.ErrNoRows {
            c.JSON(http.StatusNotFound, gin.H{"error": "Visit not found"})
        } else {
            c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error: " + err.Error()})
        }
        return
    }

    if !visitOwner.Valid || visitOwner.Int64 != int64(1) { // Assuming userID "1" corresponds to int64(1)
        c.JSON(http.StatusForbidden, gin.H{"error": "Visit does not belong to user"})
        return
    }

    // Delete visit if owner matches
    err = q.DeleteVisitByID(c, int64(visitIDInt))
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete visit: " + err.Error()})
        return
    }

    c.Status(http.StatusNoContent) // 204 No Content for successful deletion
}


func initDatabase(db *sql.DB) {
    // Load the scheme.sql and execute it
    _, err := db.Exec("path/to/scheme.sql")
    if err != nil {
        fmt.Printf("Error initializing the database: %s\n", err.Error())
        os.Exit(1)
    }
    fmt.Println("Database schema applied successfully.")
}


func main() {
	r := SetupApp()
	r.Run() // listen and serve on 0.0.0.0:8080
}
