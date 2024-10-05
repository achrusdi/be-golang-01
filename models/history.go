package models

type History struct {
	ID         string `json:"id"`
	Email      string `json:"email"`
	Password   string `json:"password"`
	IsLoggedIn bool   `json:"is_logged_in"`
}
