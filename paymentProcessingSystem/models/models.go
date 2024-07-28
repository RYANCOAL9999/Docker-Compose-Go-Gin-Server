package models

import "time"

type Describle struct {
	CardNumber     string `json:"card_number"`
	ExpirationDate string `json:"expiration_date"`
	Key            string `json:"key"`
	Currency       string `json:"currency"`
	Sender         string `json:"sender"`
	Receiver       string `json:"receiver"`
	Description    string `json:"description"`
}

// table for Payment
type Payment struct {
	ID        int       `json:"id"`
	Method    string    `json:"method" binding:"required"`
	Amount    float64   `json:"amount" binding:"required"`
	Describle Describle `json:"describle" binding:"required"`
	Timestamp time.Time `json:"timestamp"`
}

type PaymentResult struct {
	ID     int    `json:"id"`
	Status string `json:"status"`
}

type PaymentRequest struct {
	Amount         float64 `json:"amount"`
	Currency       string  `json:"currency"`
	CardNumber     string  `json:"card_number"`
	ExpirationDate string  `json:"expiration_date"`
	CVV            string  `json:"cvv"`
	Description    string  `json:"description"`
}

type TransferRequest struct {
	SenderAccount   string  `json:"sender_account"`
	ReceiverAccount string  `json:"receiver_account"`
	Amount          float64 `json:"amount"`
	Currency        string  `json:"currency"`
	Description     string  `json:"description"`
}

type BlockchainPaymentRequest struct {
	SenderAddress   string  `json:"sender_address"`
	ReceiverAddress string  `json:"receiver_address"`
	Amount          float64 `json:"amount"`
	Currency        string  `json:"currency"`
	PrivateKey      string  `json:"private_key"`
}

type PaymentResponse struct {
	TransactionID string `json:"transaction_id"`
	Status        string `json:"status"`
	Confirmations int    `json:"confirmations"`
	Message       string `json:"message"`
}

// ErrorResponse represents an error response with a single error message.
type ErrorResponse struct {
	Error string `json:"error"`
}
