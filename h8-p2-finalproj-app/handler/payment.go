package handler

import (
	"errors"
	"h8-p2-finalproj-app/model"
	"h8-p2-finalproj-app/service"
	"h8-p2-finalproj-app/util"
	"net/http"
	"os"
	"strconv"

	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

type PaymentHandler struct {
	db *gorm.DB
}

func NewPaymentHandler(db *gorm.DB) *PaymentHandler {
	return &PaymentHandler{
		db: db,
	}
}

// webhook payload structure
type WebhookPayload struct {
	ID            string  `json:"id"`
	ExternalID    string  `json:"external_id"`
	PaymentMethod string  `json:"payment_method"`
	PaidAmount    float64 `json:"paid_amount"`
	Status        string  `json:"status"`
}

type PaymentResp struct {
	ID            uint    `json:"payment_id"`
	PaymentMethod string  `json:"payment_method"`
	PaidAmount    float64 `json:"paid_amount"`
}

func (ph *PaymentHandler) GetUserForPayment(p *model.Payment) *model.User {
	if p.PurchaseType == "rentals" {
		var r model.Rental
		err := ph.db.Preload("User").Where("id=?", p.PurchaseID).First(&r).Error
		if err != nil {
			return nil
		}
		return &r.User
	}
	if p.PurchaseType == "top_ups" {
		var t model.TopUp
		err := ph.db.Preload("User").Where("id=?", p.PurchaseID).First(&t).Error
		if err != nil {
			return nil
		}
		return &t.User
	}
	return nil
}

func (ph *PaymentHandler) HandlePaymentSuccess(c echo.Context) error {
	// verify webhook token
	verifToken := c.Request().Header.Get("x-callback-token")
	if verifToken == "" || verifToken != os.Getenv("XENDIT_WEBHOOK_TOKEN") {
		return util.NewAppError(http.StatusUnauthorized, "invalid webhook token", "")
	}

	// parse body
	var reqBody WebhookPayload
	err := c.Bind(&reqBody)
	if err != nil {
		return util.NewAppError(http.StatusBadRequest, "bad request", err.Error())
	}
	if reqBody.PaymentMethod == "" || reqBody.PaidAmount <= 0 {
		return util.NewAppError(http.StatusBadRequest, "payment method and paid amount cannot be empty", "")
	}

	// get corresponding payment
	paymentId, err := strconv.Atoi(reqBody.ExternalID)
	if err != nil {
		return util.NewAppError(http.StatusBadRequest, "invalid payment id", "")
	}
	var payment model.Payment
	err = ph.db.Where("id=?", paymentId).Select("*").First(&payment).Error
	if err != nil && errors.Is(err, gorm.ErrRecordNotFound) {
		return util.NewAppError(http.StatusNotFound, "payment not found", "")
	} else if err != nil {
		return util.NewAppError(http.StatusInternalServerError, "internal server error", err.Error())
	}

	// check payment not updated yet
	if payment.Status != "Unpaid" {
		return util.NewAppError(http.StatusBadRequest, "payment already updated", "")
	}

	// update payment
	payment.PaymentMethod = reqBody.PaymentMethod
	payment.TotalPayment = reqBody.PaidAmount
	if reqBody.Status == "PAID" {
		payment.Status = "Completed"
	} else {
		payment.Status = reqBody.Status
	}
	err = ph.db.Save(&payment).Error
	if err != nil {
		return util.NewAppError(http.StatusInternalServerError, "internal server error", err.Error())
	}

	// send email is successful
	// get user
	if user := ph.GetUserForPayment(&payment); user != nil {
		if payment.Status == "Completed" {
			service.SendMail(user.Email, "Payment received!",
				"<h1>We have received your payment, enjoy your drive!</h1>",
				c.Logger(),
			)
		}
	}

	return c.JSON(http.StatusOK, map[string]any{
		"message": "payment updated successfully",
		"payment": PaymentResp{
			ID:            payment.ID,
			PaymentMethod: payment.PaymentMethod,
			PaidAmount:    payment.TotalPayment,
		},
	})
}
