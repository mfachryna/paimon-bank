package entity

type Balance struct {
	Balance  string `json:"balance"`
	Currency string `json:"currency"`
	UserId   string `json:"-"`
}
