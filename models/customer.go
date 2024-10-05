package models

type Customer struct {
	ID         string  `json:"id"`
	Name       string  `json:"name"`
	Email      string  `json:"email"`
	Password   string  `json:"password"`
	IsLoggedIn bool    `json:"is_logged_in"`
	Token      string  `json:"token"`
	Balance    float64 `json:"balance"`
}
