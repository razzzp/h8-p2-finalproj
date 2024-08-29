package model

import (
	"time"

	"gorm.io/gorm"
)

type Rental struct {
	gorm.Model
	UserID     uint
	User       User
	CarID      uint
	Car        Car
	StartDate  time.Time
	EndDate    time.Time
	TotalPrice float64
	Payment    Payment `gorm:"polymorphic:Purchase"`
}
