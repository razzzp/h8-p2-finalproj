package handler

import (
	"h8-p2-finalproj-app/model"
	"h8-p2-finalproj-app/util"
	"net/http"
	"strconv"
	"time"

	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

type CarHandler struct {
	db *gorm.DB
}

func NewCarHandler(db *gorm.DB) UserHandler {
	return UserHandler{
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

func (ch *CarHandler) GetQueryParams(c echo.Context) (*GetCarsQueryParams, error) {
	param := GetCarsQueryParams{}
	if startDate := c.QueryParam("startDate"); startDate != "" {
		// parse date
		date, err := time.Parse(time.DateOnly, startDate)
		if err != nil {
			return nil, util.NewAppError(http.StatusBadRequest, "invalid start date", "")
		}
		param.StartDate = &date
	}
	if endDate := c.QueryParam("endDate"); endDate != "" {
		// parse date
		date, err := time.Parse(time.DateOnly, endDate)
		if err != nil {
			return nil, util.NewAppError(http.StatusBadRequest, "invalid end date", "")
		}
		param.EndDate = &date
	}
	if seats := c.QueryParam("seats"); seats != "" {
		// parse date
		seats, err := strconv.Atoi(seats)
		if err != nil || seats < 1 {
			return nil, util.NewAppError(http.StatusBadRequest, "invalid number of seats", "")
		}
		param.Seats = &seats
	}
	return &param, nil
}

func (ch *CarHandler) GetCars(params *GetCarsQueryParams) ([]*model.Car, error) {
	q := ch.db.Preload("Car").Model(&model.Rental{}).Select("")
	if params.StartDate != nil {
		q = q.
			Where("start_date > ? AND start_date < ?")
	}
	if params.EndDate != nil {
		q = q.
			Or("end_date > ? AND start_date < ?")
	}
	q = q.Group("car_id")
	return nil, nil
}

func (ch *CarHandler) HandleGetCars(c echo.Context) error {
	params, err := ch.GetQueryParams(c)
	if err != nil {
		return err
	}
	availCars, err := ch.GetCars(params)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, availCars)
}
