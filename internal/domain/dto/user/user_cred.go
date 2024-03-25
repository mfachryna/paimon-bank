package dto

type UserCredential struct {
	Name        string `json:"name"`
	Email       string `json:"email"`
	Phone       string `json:"phone"`
	AccessToken string `json:"accessToken"`
}
