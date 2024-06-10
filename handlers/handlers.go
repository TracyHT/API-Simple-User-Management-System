package handlers

import (
	"database/sql"
	"net/http"

	"github.com/gin-gonic/gin"

	"main.go/database"
	models "main.go/model"
)

func CreateUser(c *gin.Context) {
	var user models.User
	if err := c.BindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if user.Username == "" || user.Firstname == "" || user.Lastname == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "username, firstname, and lastname are required"})
		return
	}

	stmt, err := database.DB.Prepare("INSERT INTO users(username, firstname, lastname, email, avatar, phone, date_of_birth, address_country, address_city, address_street_name, address_street_address) VALUES(?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)")
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

func GetUser(c *gin.Context) {
	id := c.Param("id")
	var user models.User
	err := database.DB.QueryRow("SELECT id, username, firstname, lastname, email, avatar, phone, date_of_birth, address_country, address_city, address_street_name, address_street_address FROM users WHERE id = ?", id).Scan(&user.ID, &user.Username, &user.Firstname, &user.Lastname, &user.Email, &user.Avatar, &user.Phone, &user.DateOfBirth, &user.AddressCountry, &user.AddressCity, &user.AddressStreetName, &user.AddressStreetAddr)
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

func UpdateUser(c *gin.Context) {
	id := c.Param("id")
	var user models.User
	if err := c.BindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if user.Username == "" || user.Firstname == "" || user.Lastname == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "username, firstname, and lastname are required"})
		return
	}

	stmt, err := database.DB.Prepare("UPDATE users SET username = ?, firstname = ?, lastname = ?, email = ?, avatar = ?, phone = ?, date_of_birth = ?, address_country = ?, address_city = ?, address_street_name = ?, address_street_address = ? WHERE id = ?")
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

func DeleteUser(c *gin.Context) {
	id := c.Param("id")
	stmt, err := database.DB.Prepare("DELETE FROM users WHERE id = ?")
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

func ListUsers(c *gin.Context) {
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
	rows, err := database.DB.Query(query, args...)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer rows.Close()

	var users []models.User
	for rows.Next() {
		var user models.User
		if err := rows.Scan(&user.ID, &user.Username, &user.Firstname, &user.Lastname, &user.Email, &user.Avatar, &user.Phone, &user.DateOfBirth, &user.AddressCountry, &user.AddressCity, &user.AddressStreetName, &user.AddressStreetAddr); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		users = append(users, user)
	}
	c.JSON(http.StatusOK, users)
}
