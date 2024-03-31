package entity

type Balance struct {
	Balance  int64  `json:"balance"`
	Currency string `json:"currency"`
	UserId   string `json:"-"`
}
