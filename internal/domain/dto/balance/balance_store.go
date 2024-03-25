package dto

type BalanceAdd struct {
	BankNumber string `json:"senderBankAccountNumber" validate:"required,min=5,max=30"`
	BankName   string `json:"senderBankName" validate:"required,min=5,max=30"`
	Balance    int64  `json:"addedBalance" validate:"required,min=1"`
	Currency   string `json:"currency" validate:"required,iso4217"`
	Receipt    string `json:"transferProofImg" validate:"required,url"`
}
type BalanceTransaction struct {
	BankNumber string `json:"recipientBankAccountNumber" validate:"required,min=5,max=30"`
	BankName   string `json:"recipientBankName" validate:"required,min=5,max=30"`
	Balance    int64  `json:"balances" validate:"required,min=1"`
	Currency   string `json:"fromCurrency" validate:"required,iso4217"`
	Receipt    string `json:"transferProofImg" validate:"required,url"`
}
