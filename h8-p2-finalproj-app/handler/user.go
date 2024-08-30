package handler

import (
	"errors"
	"fmt"
	"h8-p2-finalproj-app/auth"
	"h8-p2-finalproj-app/model"
	"h8-p2-finalproj-app/service"
	"h8-p2-finalproj-app/util"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type UserHandler struct {
	db *gorm.DB
	is *service.InvoiceService
}

func NewUserHandler(db *gorm.DB, is *service.InvoiceService) UserHandler {
	return UserHandler{
		db: db,
		is: is,
	}
}

type RegisterReqData struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type RegisterRespData struct {
	ID    uint   `json:"id,omitempty"`
	Name  string `json:"name"`
	Email string `json:"email"`
}

type LoginReqData struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type LoginRespData struct {
	Token string `json:"token"`
}

func (uh *UserHandler) validateRegisterUserData(ud *RegisterReqData) error {
	// check name
	ud.Name = strings.TrimSpace(ud.Name)
	if ud.Name == "" {
		return util.NewAppError(http.StatusBadRequest, "full name cannot be empty", "")
	}
	// check email
	ud.Email = strings.TrimSpace(ud.Email)
	if ud.Email == "" {
		return util.NewAppError(http.StatusBadRequest, "email cannot be empty", "")
	}
	// check dupe
	var count int64
	err := uh.db.Model(&model.User{}).Where("email=?", ud.Email).Count(&count).Error
	if err != nil {
		return util.NewAppError(http.StatusInternalServerError, "internal server error", err.Error())
	}
	if count != 0 {
		// email already registered
		return util.NewAppError(http.StatusBadRequest, "email already registered", "")
	}

	// check pass cannot be empty and less than 4 chars
	if ud.Password == "" {
		return util.NewAppError(http.StatusBadRequest, "password cannot be empty", "")
	}
	if len(ud.Password) < 4 {
		return util.NewAppError(http.StatusBadRequest, "password must be at least 4 characters", "")
	}

	return nil
}

//	@Summary	Registers a user
//	@Tags		users
//	@Accept		json
//	@Param		UserData	body	handler.RegisterReqData	true	"Data of user to register"
//	@Produce	json
//	@Success	201	{object}	handler.RegisterRespData
//	@Failure	400	{object}	util.AppError
//	@Failure	500	{object}	util.AppError
//	@Router		/users/register [post]
func (uh *UserHandler) HandleRegisterUser(c echo.Context) error {
	// get body
	var userData RegisterReqData
	err := c.Bind(&userData)
	if err != nil {
		return util.NewAppError(http.StatusInternalServerError, "failed to parse request body", err.Error())
	}

	// validate body
	err = uh.validateRegisterUserData(&userData)
	if err != nil {
		return err
	}

	// gen password
	passHash, err := bcrypt.GenerateFromPassword([]byte(userData.Password), 12)
	if err != nil {
		return util.NewAppError(http.StatusInternalServerError, "internal server error", err.Error())
	}
	// new entity
	newUser := model.User{
		Name:     userData.Name,
		Email:    userData.Email,
		Password: string(passHash),
	}
	err = uh.db.Create(&newUser).Error
	if err != nil {
		return util.NewAppError(http.StatusInternalServerError, "internal server error", err.Error())
	}

	// create response data
	resp := RegisterRespData{
		ID:    newUser.ID,
		Name:  newUser.Name,
		Email: newUser.Email,
	}

	err = service.SendMail(userData.Email, "Welcome to car rental app!", fmt.Sprintf("<h1>Hello %s, thank you for registering with us!</h1>", newUser.Name))
	if err != nil {
		c.Logger().Errorf("failed to send email notif: %s", err.Error())
	}

	return c.JSON(http.StatusCreated, &util.ResponseData{
		Message: "User successfully registered",
		Data:    &resp,
	})
}

func (uh *UserHandler) validateLoginData(ud *RegisterReqData) error {
	if ud.Email == "" {
		return util.NewAppError(http.StatusInternalServerError, "email cannot be empty", "")
	}
	if ud.Password == "" {
		return util.NewAppError(http.StatusInternalServerError, "password cannot be empty", "")
	}
	return nil
}

//	@Summary	Login
//	@Tags		users
//	@Accept		json
//	@Param		EmailPassword	body	handler.LoginReqData	true	"Email and password"
//	@Produce	json
//	@Success	200	{object}	handler.LoginRespData
//	@Failure	404	{object}	util.AppError
//	@Failure	400	{object}	util.AppError
//	@Failure	500	{object}	util.AppError
//	@Router		/users/login [post]
func (uh *UserHandler) HandleLoginUser(c echo.Context) error {

	// parse body
	var loginData RegisterReqData
	err := c.Bind(&loginData)
	if err != nil {
		return util.NewAppError(http.StatusInternalServerError, "failed to parse request body", err.Error())
	}

	// validate data
	err = uh.validateLoginData(&loginData)
	if err != nil {
		return err
	}
	// get existing user
	var existingUser model.User
	err = uh.db.Where("email =?", loginData.Email).First(&existingUser).Error
	if err != nil && errors.Is(err, gorm.ErrRecordNotFound) {
		// email not registed
		return util.NewAppError(http.StatusNotFound, "email not registered", "")
	} else if err != nil {
		// other error
		return util.NewAppError(http.StatusInternalServerError, "internal server error", err.Error())
	}
	fmt.Println(existingUser)
	// compare password
	err = bcrypt.CompareHashAndPassword([]byte(existingUser.Password), []byte(loginData.Password))
	if err != nil {
		fmt.Println(err)
		return util.NewAppError(http.StatusBadRequest, "incorrect password", "")
	}

	// Set custom claims
	claims := &auth.JwtAppClaims{
		UserID: existingUser.ID,
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   existingUser.Email,
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 5)),
		},
	}

	// Create token with claims
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Generate encoded token and send it as response.
	t, err := token.SignedString([]byte(os.Getenv("JWT_KEY")))
	if err != nil {
		return util.NewAppError(http.StatusInternalServerError, "internal server error", err.Error())
	}

	return c.JSON(http.StatusOK, map[string]string{
		"token": t,
	})
}

type UserProfile struct {
	UserID  uint    `json:"user_id"`
	Name    string  `json:"name"`
	Email   string  `json:"email"`
	Deposit float64 `json:"deposit"`
}

//	@Summary	User profile
//	@Tags		users
//	@Produce	json
//	@Success	200	{object}	handler.UserProfile
//	@Failure	401	{object}	util.AppError
//	@Failure	500	{object}	util.AppError
//	@Router		/users/profile [get]
func (uh *UserHandler) HandleUserProfile(c echo.Context) error {
	// get user from context
	user, err := util.GetUserFromContext(c, uh.db)
	if err != nil {
		return err
	}
	resp := UserProfile{
		UserID:  user.ID,
		Name:    user.Name,
		Email:   user.Email,
		Deposit: user.Deposit,
	}

	return c.JSON(http.StatusOK, resp)
}

type TopUpReq struct {
	Amount float64 `json:"amount"`
}
type TopUpResp struct {
	TopUpID       uint    `json:"top_up_id"`
	Amount        float64 `json:"amount"`
	PaymentID     uint    `json:"payment_id"`
	PaymentStatus string  `json:"payment_status"`
	PaymentUrl    string  `json:"payment_url"`
}

//	@Summary	top up user deposit
//	@Tags		users
//	@Produce	json
//	@Accept		json
//	@Param		Amount	body		handler.TopUpReq	true	"Amount to top up"
//	@Success	200		{object}	handler.TopUpResp
//	@Failure	401		{object}	util.AppError
//	@Failure	500		{object}	util.AppError
//	@Router		/users/topup [post]
func (uh *UserHandler) HandlePostTopUp(c echo.Context) error {
	// get user from context
	user, err := util.GetUserFromContext(c, uh.db)
	if err != nil {
		return err
	}
	c.Logger().Printf("User found: %s", user.Email)

	// parse and validate req body
	var reqBody TopUpReq
	err = c.Bind(&reqBody)
	if err != nil {
		return util.NewAppError(http.StatusBadRequest, "bad request", err.Error())
	}

	if reqBody.Amount <= 0 {
		return util.NewAppError(http.StatusBadRequest, "amount cannot be 0", "")
	}

	newTopUp := model.TopUp{
		UserID: user.ID,
		Amount: reqBody.Amount,
	}

	newPayment := model.Payment{
		PaymentUrl:    "",
		Status:        "Unpaid",
		PaymentMethod: "",
		TotalPayment:  0,
	}
	newTopUp.Payment = newPayment
	newTopUp.User = *user

	err = uh.db.Create(&newTopUp).Error
	if err != nil {
		return util.NewAppError(http.StatusInternalServerError, "internal server error", err.Error())
	}

	// generate invoice url
	url, err := uh.is.GenerateInvoice(
		newTopUp.Payment.ID,
		newTopUp.Amount,
		fmt.Sprintf("Top up IDR %.0f", newTopUp.Amount),
		user.Email,
	)
	if err != nil {
		return util.NewAppError(http.StatusInternalServerError, "internal server error", err.Error())
	}
	// update payment url
	// use one assigned to rental !!!!
	newTopUp.Payment.PaymentUrl = url
	err = uh.db.Save(&newTopUp.Payment).Error
	if err != nil {
		return util.NewAppError(http.StatusInternalServerError, "internal server error", err.Error())
	}

	resp := TopUpResp{
		TopUpID:       newTopUp.ID,
		Amount:        newTopUp.Amount,
		PaymentID:     newTopUp.Payment.ID,
		PaymentStatus: newTopUp.Payment.Status,
		PaymentUrl:    url,
	}

	return c.JSON(http.StatusCreated, resp)
}
