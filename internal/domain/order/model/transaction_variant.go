package model

type TransactionVariant struct {
	ID            int    `json:"id"`
	Value         string `json:"value"`
	TransactionID int    `json:"transactionId"`
}
