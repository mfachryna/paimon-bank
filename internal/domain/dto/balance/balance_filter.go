package dto

type BalanceFilter struct {
	Limit  int64 `json:"limit" validate:"omitempty,numeric,min=0" schema:"limit"`
	Offset int64 `json:"offset" validate:"omitempty,numeric,min=0" schema:"offset"`
}
