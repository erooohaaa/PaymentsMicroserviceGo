package domain

import "errors"

var ErrNotFound = errors.New("payment not found")

type PaymentStatus = string

const (
	StatusAuthorized PaymentStatus = "Authorized"
	StatusDeclined   PaymentStatus = "Declined"
)


const MaxPaymentAmount int64 = 100000

type Payment struct {
	ID            string
	OrderID       string
	TransactionID string
	Amount        int64
	Status        PaymentStatus
}
