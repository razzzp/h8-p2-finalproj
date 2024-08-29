package handler

import (
	"h8-p2-finalproj-app/service"
	"h8-p2-finalproj-app/util"
	"net/http"
	"strconv"
	"time"

	"github.com/labstack/echo/v4"
)

type CarHandler struct {
	cs *service.CarService
}

func NewCarHandler(cs *service.CarService) CarHandler {
	return CarHandler{
		cs: cs,
	}
}

func (ch *CarHandler) GetQueryParams(c echo.Context) (*service.GetCarsQueryParams, error) {
	param := service.GetCarsQueryParams{}
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

type GetCarsRespItem struct {
	CardID             uint   `json:"car_id"`
	Manufacturer       string `json:"manufacturer"`
	CarModel           string `json:"model"`
	Seats              uint   `json:"seats"`
	NumOfCarsAvailable uint   `json:"num_of_cars_available"`
}

func (ch *CarHandler) HandleGetCars(c echo.Context) error {
	params, err := ch.GetQueryParams(c)
	if err != nil {
		return err
	}
	availCars, err := ch.cs.GetCarsWithRentals(params)
	if err != nil {
		return err
	}
	// convert to brief with available cars within that date range
	resp := []GetCarsRespItem{}
	for _, ac := range availCars {
		if ac.Stock <= ac.NumOfRentals {
			// skip cars that are fully rented out
			continue
		}
		resp = append(resp, GetCarsRespItem{
			CardID:             ac.ID,
			Manufacturer:       ac.Manufacturer,
			CarModel:           ac.CarModel,
			Seats:              ac.Seats,
			NumOfCarsAvailable: ac.Stock - ac.NumOfRentals,
		})
	}
	return c.JSON(http.StatusOK, resp)
}
