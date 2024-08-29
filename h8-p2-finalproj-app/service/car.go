package service

import (
	"encoding/json"
	"fmt"
	"h8-p2-finalproj-app/model"
	"time"

	"gorm.io/gorm"
)

type CarService struct {
	db *gorm.DB
}

func NewCarService(db *gorm.DB) *CarService {
	return &CarService{
		db: db,
	}
}

type GetCarsQueryParams struct {
	StartDate    *time.Time
	EndDate      *time.Time
	WheelDrive   *int
	Type         []string
	Seats        *int
	Transmission *string
	Manufacturer []string
	CarModel     []string
	Year         *uint
}

func (cs *CarService) CountRentalsPerCarQuery(params *GetCarsQueryParams) *gorm.DB {
	q := cs.db.Model(&model.Rental{}).Select("car_id, COUNT(rentals.id) AS num_of_rentals")
	if params.StartDate != nil && params.EndDate != nil {
		q = q.
			Where("start_date >= ? AND start_date <= ?", *params.StartDate, *params.EndDate).
			Or("end_date >= ? AND end_date <= ?", *params.StartDate, *params.EndDate).
			Or("start_date <= ? AND end_date >= ?", *params.StartDate, *params.EndDate)

	}
	q = q.Group("car_id")
	return q
}

func (cs *CarService) CarsWithRentalQuery(params *GetCarsQueryParams) *gorm.DB {
	countRentalQ := cs.CountRentalsPerCarQuery(params)
	q := cs.db.
		Model(&model.Car{}).
		Joins("left join (?) q on cars.id = q.car_id", countRentalQ)
	if params.Seats != nil {
		q = q.Where("cars.seats >= ?", *params.Seats)
	}
	q = q.Select("*")
	return q
}

// struct for easily getting query
type AvailableCarData struct {
	model.Car
	CarID        uint
	NumOfRentals uint
}

func (cs *CarService) GetCarsWithRentals(params *GetCarsQueryParams) ([]AvailableCarData, error) {

	q := cs.CarsWithRentalQuery(&GetCarsQueryParams{})

	var result []AvailableCarData

	tx := q.Find(&result)
	if tx.Error != nil {
		// fmt.Println(tx.Error)
		return nil, tx.Error
	}
	bytes, _ := json.Marshal(result)
	fmt.Println(string(bytes))
	return result, nil
}
