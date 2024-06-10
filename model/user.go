package models

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
