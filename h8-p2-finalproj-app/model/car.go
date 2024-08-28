package model

import "gorm.io/gorm"

type Car struct {
	gorm.Model
	WheelDrive   uint
	Type         string
	Seats        uint
	Transmission string
	Manufacturer string
	CarModel     string
	Year         uint
	Stock        uint
	RatePerDay   float64
}
