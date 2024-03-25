package dto

import "time"

type UserCreate struct {
	Email     string    `json:"email"`
	Name      string    `json:"name" validate:"required,min=5,max=50"`
	Password  string    `json:"password" validate:"required,min=5,max=15"`
	CreatedAt time.Time `json:"createdAt"`
}
