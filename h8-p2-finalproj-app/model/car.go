package model

import (
	"fmt"

	"gorm.io/gorm"
)

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

func (c *Car) GetCarName() string {
	return fmt.Sprintf("%s %s", c.Manufacturer, c.CarModel)
}
