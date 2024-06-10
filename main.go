package main

import (
	"fmt"
	"os"

	"github.com/gin-gonic/gin"
	_ "github.com/mattn/go-sqlite3"

	"main.go/database"
	"main.go/routes"
)

func main() {
	// Print the value of CGO_ENABLED
	fmt.Println("CGO_ENABLED:", os.Getenv("CGO_ENABLED"))

	// Initialize the database
	db, err := database.InitDB("./user_management.db")
	if err != nil {
		fmt.Println("Error opening database:", err)
		return
	}
	defer db.Close()

	// Set up the Gin router
	r := gin.Default()
	routes.SetupRoutes(r)

	// Run the server
	r.Run(":8080")
}
