package models

type User struct {
	ID       int    `json:"id"`
	Username string `json:"username"`
	Password string `json:"-"`
	IsAdmin  bool   `json:"is_admin"`
}
