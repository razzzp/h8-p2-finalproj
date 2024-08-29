package handler

import (
	"errors"
	"h8-p2-finalproj-app/model"
	"h8-p2-finalproj-app/service"
	"h8-p2-finalproj-app/util"
	"math"
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

type RentalHandler struct {
	db *gorm.DB
	cs *service.CarService
	is *service.InvoiceService
}

func NewRentalHandler(
	db *gorm.DB,
	cs *service.CarService,
	is *service.InvoiceService) RentalHandler {
	return RentalHandler{
		db: db,
		cs: cs,
		is: is,
	}
}

type PostRentalsReq struct {
	CarID     uint   `json:"car_id"`
	StartDate string `json:"start_date"`
	EndDate   string `json:"end_date"`
}

type PostRentalData struct {
	CarID     uint      `json:"car_id"`
	StartDate time.Time `json:"start_date"`
	EndDate   time.Time `json:"end_date"`
}

func (rh *RentalHandler) validatePostRentalReqData(prr *PostRentalsReq) (*PostRentalData, error) {
	startDate, err := time.Parse(time.DateOnly, prr.StartDate)
	if err != nil {
		return nil, util.NewAppError(http.StatusBadRequest, "invalid end date", "")
	}

	endDate, err := time.Parse(time.DateOnly, prr.EndDate)
	if err != nil {
		return nil, util.NewAppError(http.StatusBadRequest, "invalid start date", "")
	}
	if endDate.Before(startDate) {
		return nil, util.NewAppError(http.StatusBadRequest, "end date cannot be before start date", "")
	}
	return &PostRentalData{
		CarID:     prr.CarID,
		StartDate: startDate,
		EndDate:   endDate,
	}, nil
}

type PostRentalResp struct {
	RentalID      uint      `json:"rental_id"`
	CarID         uint      `json:"car_id"`
	StartDate     time.Time `json:"start_date"`
	EndDate       time.Time `json:"end_date"`
	TotalPrice    float64   `json:"total_price"`
	PaymentStatus string    `json:"payment_status"`
	PaymentUrl    string    `json:"payment_url"`
}

func (rh *RentalHandler) HandlePostRentals(c echo.Context) error {
	// get user from context
	user, err := util.GetUserFromContext(c, rh.db)
	if err != nil {
		return err
	}

	// parse and validate req body
	var reqBody PostRentalsReq
	err = c.Bind(&reqBody)
	if err != nil {
		return util.NewAppError(http.StatusBadRequest, "bad request", err.Error())
	}
	rentalData, err := rh.validatePostRentalReqData(&reqBody)
	if err != nil {
		return err
	}

	// get car details
	var car model.Car
	err = rh.db.Where("id=?", rentalData.CarID).First(&car).Error
	if err != nil && errors.Is(err, gorm.ErrRecordNotFound) {
		return util.NewAppError(http.StatusNotFound, "car not found", "")
	} else if err != nil {
		return util.NewAppError(http.StatusInternalServerError, "internal server error", err.Error())
	}

	// check car available
	isAvail, err := rh.cs.IsCarAvailable(car.ID, rentalData.StartDate, rentalData.EndDate)
	if err != nil {
		return util.NewAppError(http.StatusInternalServerError, "internal server error", err.Error())
	}
	if !isAvail {
		return util.NewAppError(http.StatusBadRequest, "car is not available", "")
	}

	// calculate total price, and create rental
	days := math.Ceil(rentalData.EndDate.Sub(rentalData.StartDate).Hours() / 24)
	newRental := model.Rental{
		UserID:     user.ID,
		CarID:      car.ID,
		StartDate:  rentalData.StartDate,
		EndDate:    rentalData.EndDate,
		TotalPrice: days * car.RatePerDay,
	}
	newPayment := model.Payment{
		PaymentUrl:    "https://example.com",
		Status:        "Unpaid",
		PaymentMethod: "",
		TotalPayment:  0,
	}
	newRental.Payment = newPayment
	newRental.User = *user

	err = rh.db.Create(&newRental).Error
	if err != nil {
		return util.NewAppError(http.StatusInternalServerError, "internal server error", err.Error())
	}

	// generate invoice url
	url, err := rh.is.GenerateInvoice(&newRental)
	if err != nil {
		return util.NewAppError(http.StatusInternalServerError, "internal server error", err.Error())
	}

	resp := PostRentalResp{
		RentalID:      newRental.ID,
		StartDate:     newRental.StartDate,
		EndDate:       newRental.EndDate,
		TotalPrice:    newRental.TotalPrice,
		PaymentStatus: newPayment.Status,
		PaymentUrl:    url,
	}

	return c.JSON(http.StatusCreated, resp)
}
