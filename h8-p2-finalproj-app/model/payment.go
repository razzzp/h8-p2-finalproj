package model

import "gorm.io/gorm"

type Payment struct {
	gorm.Model
	PurchaseID    int    `gorm:"column:purchase_id;not null"`
	PurchaseType  string `gorm:"column:purchase_type;not null"`
	PaymentUrl    string
	Status        string
	PaymentMethod string
	TotalPayment  float64
}
