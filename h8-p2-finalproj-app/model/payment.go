package model

import "gorm.io/gorm"

type Payment struct {
	gorm.Model
	PurchaseID    int
	PurchaseType  string
	PaymentUrl    string
	Status        string
	PaymentMethod string
	TotalPayment  float64
}
