package model

import (
	"gorm.io/gorm"
)

type TopUp struct {
	gorm.Model
	UserID  uint
	User    User
	Amount  float64
	Payment Payment `gorm:"polymorphicType:PurchaseType;polymorphicId:PurchaseID"`
}
