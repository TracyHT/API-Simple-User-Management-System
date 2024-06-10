package main

import (
	"database/sql"
	"fmt"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	_ "github.com/mattn/go-sqlite3"
)

type User struct {
	ID                int    `json:"id"`
	Username          string `json:"username"`
	Firstname         string `json:"firstname"`
	Lastname          string `json:"lastname"`
	Email             string `json:"email"`
	Avatar            string `json:"avatar"`
	Phone             string `json:"phone"`
	DateOfBirth       string `json:"date_of_birth"`
	AddressCountry    string `json:"address_country"`
	AddressCity       string `json:"address_city"`
	AddressStreetName string `json:"address_street_name"`
	AddressStreetAddr string `json:"address_street_address"`
}

var db *sql.DB

func main() {
	// Print the value of CGO_ENABLED
	fmt.Println("CGO_ENABLED:", os.Getenv("CGO_ENABLED"))

	var err error
	db, err = sql.Open("sqlite3", "./user_management.db")
	if err != nil {
		fmt.Println("Error opening database:", err)
		return
	}
	defer db.Close()

	// Create 'users' table
	_, err = db.Exec(`
        CREATE TABLE IF NOT EXISTS users (
    	id INTEGER PRIMARY KEY AUTOINCREMENT,
    	username TEXT NOT NULL UNIQUE,
    	firstname TEXT NOT NULL,
   		lastname TEXT NOT NULL,
    	email TEXT UNIQUE,
    	avatar TEXT,
    	phone TEXT,
    	date_of_birth DATE CHECK(date_of_birth IS NULL OR date_of_birth <= date('now')), -- Example check constraint for date_of_birth
    	address_country TEXT,
    	address_city TEXT,
    	address_street_name TEXT,
    	address_street_address TEXT
);
    `)
	if err != nil {
		fmt.Println("Error creating table:", err)
		return
	} else {
		fmt.Println("Table created successfully")
	}

	r := gin.Default()
	r.POST("/users", createUser)
	r.GET("/users/:id", getUser)
	r.PUT("/users/:id", updateUser)
	r.DELETE("/users/:id", deleteUser)
	r.GET("/users", listUsers)

	r.Run(":8080")
}

func createUser(c *gin.Context) {
	var user User
	if err := c.BindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if user.Username == "" || user.Firstname == "" || user.Lastname == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "username, firstname, and lastname are required"})
		return
	}

	stmt, err := db.Prepare("INSERT INTO users(username, firstname, lastname, email, avatar, phone, date_of_birth, address_country, address_city, address_street_name, address_street_address) VALUES(?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	res, err := stmt.Exec(user.Username, user.Firstname, user.Lastname, user.Email, user.Avatar, user.Phone, user.DateOfBirth, user.AddressCountry, user.AddressCity, user.AddressStreetName, user.AddressStreetAddr)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	id, _ := res.LastInsertId()
	user.ID = int(id)
	c.JSON(http.StatusOK, user)
}

func getUser(c *gin.Context) {
	id := c.Param("id")
	var user User
	err := db.QueryRow("SELECT id, username, firstname, lastname, email, avatar, phone, date_of_birth, address_country, address_city, address_street_name, address_street_address FROM users WHERE id = ?", id).Scan(&user.ID, &user.Username, &user.Firstname, &user.Lastname, &user.Email, &user.Avatar, &user.Phone, &user.DateOfBirth, &user.AddressCountry, &user.AddressCity, &user.AddressStreetName, &user.AddressStreetAddr)
	if err != nil {
		if err == sql.ErrNoRows {
			c.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}
	c.JSON(http.StatusOK, user)
}

func updateUser(c *gin.Context) {
	id := c.Param("id")
	var user User
	if err := c.BindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if user.Username == "" || user.Firstname == "" || user.Lastname == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "username, firstname, and lastname are required"})
		return
	}

	stmt, err := db.Prepare("UPDATE users SET username = ?, firstname = ?, lastname = ?, email = ?, avatar = ?, phone = ?, date_of_birth = ?, address_country = ?, address_city = ?, address_street_name = ?, address_street_address = ? WHERE id = ?")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	_, err = stmt.Exec(user.Username, user.Firstname, user.Lastname, user.Email, user.Avatar, user.Phone, user.DateOfBirth, user.AddressCountry, user.AddressCity, user.AddressStreetName, user.AddressStreetAddr, id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, user)
}

func deleteUser(c *gin.Context) {
	id := c.Param("id")
	stmt, err := db.Prepare("DELETE FROM users WHERE id = ?")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	_, err = stmt.Exec(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "user deleted"})
}

func listUsers(c *gin.Context) {
	// Extract query parameters for filtering and sorting
	username := c.Query("username")
	email := c.Query("email")
	firstname := c.Query("firstname")
	lastname := c.Query("lastname")
	country := c.Query("address_country")
	sortBy := c.DefaultQuery("sort_by", "id")        // Default sorting by ID
	sortOrder := c.DefaultQuery("sort_order", "asc") // Default sorting order ascending

	// Construct SQL query based on filtering and sorting criteria
	query := "SELECT id, username, firstname, lastname, email, avatar, phone, date_of_birth, address_country, address_city, address_street_name, address_street_address FROM users WHERE 1=1"
	var args []interface{}

	if username != "" {
		query += " AND username = ?"
		args = append(args, username)
	}
	if email != "" {
		query += " AND email = ?"
		args = append(args, email)
	}
	if firstname != "" {
		query += " AND firstname = ?"
		args = append(args, firstname)
	}
	if lastname != "" {
		query += " AND lastname = ?"
		args = append(args, lastname)
	}
	if country != "" {
		query += " AND address_country = ?"
		args = append(args, country)
	}

	// Add sorting criteria
	query += " ORDER BY " + sortBy + " " + sortOrder

	// Execute query with filtering and sorting criteria
	rows, err := db.Query(query, args...)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer rows.Close()

	var users []User
	for rows.Next() {
		var user User
		if err := rows.Scan(&user.ID, &user.Username, &user.Firstname, &user.Lastname, &user.Email, &user.Avatar, &user.Phone, &user.DateOfBirth, &user.AddressCountry, &user.AddressCity, &user.AddressStreetName, &user.AddressStreetAddr); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		users = append(users, user)
	}
	c.JSON(http.StatusOK, users)
}
