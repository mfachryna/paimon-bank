package entity

type BalanceHistory struct {
	ID        string        `json:"transactionId"`
	Balance   int64         `json:"balance"`
	Currency  string        `json:"currency"`
	Receipt   string        `json:"transaferProofImg"`
	CreatedAt string        `json:"createdAt"`
	Source    BalanceSource `json:"source"`
	UserId    string        `json:"-"`
}
type BalanceSource struct {
	BankNumber string `json:"bankAccountNumber"`
	BankName   string `json:"bankName"`
}
