package entity

type User struct {
	ID          string `json:"userId"`
	Email       string `json:"-"`
	Password    string `json:"-"`
	Name        string `json:"name"`
	FriendCount int64  `json:"friendCount"`
	CreatedAt   string `json:"createdAt"`
}
