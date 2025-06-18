package models

import "time"

type Payment struct {
	ID 			 string `json:"_id" bson:"_id, omitempty"`
	PaymentID    string    `json:"payment_id" bson:"payment_id"`
	OrderID      string    `json:"order_id" bson:"order_id"`
	UserID       string    `json:"user_id" bson:"user_id"`
	Amount       float64   `json:"amount" bson:"amount"`
	PaymentMethod string   `json:"payment_method" bson:"payment_method"`
	Status       string    `json:"status" bson:"status"`
	TransactionID string   `json:"transaction_id" bson:"transaction_id"`
	CreatedAt    time.Time `json:"created_at" bson:"created_at"`
	UpdatedAt    time.Time `json:"updated_at" bson:"updated_at"`
}
