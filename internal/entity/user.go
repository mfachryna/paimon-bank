package entity

type User struct {
	ID       string `json:"userId"`
	Email    string `json:"-"`
	Password string `json:"-"`
	Name     string `json:"name"`
}
