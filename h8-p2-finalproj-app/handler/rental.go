package handler

import (
	"errors"
	"fmt"
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
	PaymentID     uint      `json:"payment_id"`
	PaymentStatus string    `json:"payment_status"`
	PaymentUrl    string    `json:"payment_url"`
}

func (rh *RentalHandler) GenerateInvoiceDesc(rental *model.Rental) string {
	return fmt.Sprintf(
		"Renting %s %s, from: %s to %s",
		rental.Car.Manufacturer,
		rental.Car.CarModel,
		rental.StartDate.Format(time.DateOnly),
		rental.EndDate.Format(time.DateOnly),
	)
}

func (rh *RentalHandler) HandlePostRentals(c echo.Context) error {
	// get user from context
	user, err := util.GetUserFromContext(c, rh.db)
	if err != nil {
		return err
	}
	c.Logger().Printf("User found: %s", user.Email)

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
		PaymentUrl:    "",
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
	url, err := rh.is.GenerateInvoice(
		newRental.Payment.ID,
		newRental.TotalPrice,
		rh.GenerateInvoiceDesc(&newRental),
		user.Email,
	)
	if err != nil {
		return util.NewAppError(http.StatusInternalServerError, "internal server error", err.Error())
	}
	// update payment url
	// use one assigned to rental !!!!
	newRental.Payment.PaymentUrl = url
	err = rh.db.Save(&newRental.Payment).Error
	if err != nil {
		return util.NewAppError(http.StatusInternalServerError, "internal server error", err.Error())
	}

	// notify booking made
	err = service.SendMail(
		user.Email,
		"You've made a booking!",
		fmt.Sprintf(`
		<h1>Hello %s,</h1><br>
		<p>Thank you for making a booking, here are the details:</p>
		<p>Car: %s<br>
		Start: %s<br>
		End: %s<br>
		Total Price: IDR %.0f<br>	
		</p>
		<p>You can make payment here:<br>%s</p>
		`, user.Name,
			car.GetCarName(),
			newRental.StartDate.Format(time.DateOnly),
			newRental.EndDate.Format(time.DateOnly),
			newRental.TotalPrice,
			newPayment.PaymentUrl),
		c.Logger(),
	)
	if err != nil {
		c.Logger().Errorf("failed to send email notif: %s", err.Error())
	}

	resp := PostRentalResp{
		RentalID:      newRental.ID,
		StartDate:     newRental.StartDate,
		EndDate:       newRental.EndDate,
		TotalPrice:    newRental.TotalPrice,
		PaymentID:     newRental.Payment.ID,
		PaymentStatus: newRental.Payment.Status,
		PaymentUrl:    url,
	}

	return c.JSON(http.StatusCreated, resp)
}

type RentalRespItem struct {
	RentalID      uint   `json:"rental_id"`
	StartDate     string `json:"start_date"`
	EndDate       string `json:"end_date"`
	CarID         uint   `json:"car_id"`
	Car           string `json:"car"`
	PaymentID     uint   `json:"payment_id"`
	PaymentStatus string `json:"payment_status"`
	PaymentUrl    string `json:"payment_url,omitempty"`
}

func (rh *RentalHandler) HandleGetRentals(c echo.Context) error {
	// get user from context
	user, err := util.GetUserFromContext(c, rh.db)
	if err != nil {
		return err
	}
	c.Logger().Print(user.Email)

	// get rentals for user
	var rentals []model.Rental
	err = rh.db.Preload("Car").Preload("Payment").Where("user_id=?", user.ID).Find(&rentals).Error
	if err != nil {
		return util.NewAppError(http.StatusInternalServerError, "internal server error", err.Error())
	}

	resp := []RentalRespItem{}
	for _, r := range rentals {
		ri := RentalRespItem{
			RentalID:      r.ID,
			StartDate:     r.StartDate.Format(time.DateOnly),
			EndDate:       r.EndDate.Format(time.DateOnly),
			CarID:         r.Car.ID,
			Car:           r.Car.GetCarName(),
			PaymentID:     r.Payment.ID,
			PaymentStatus: r.Payment.Status,
		}
		if r.Payment.Status == "Unpaid" {
			ri.PaymentUrl = r.Payment.PaymentUrl
		}
		resp = append(resp, ri)
	}

	return c.JSON(http.StatusOK, resp)
}
