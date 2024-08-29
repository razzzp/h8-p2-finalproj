package handler

import (
	"errors"
	"fmt"
	"h8-p2-finalproj-app/auth"
	"h8-p2-finalproj-app/model"
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
}

func NewUserHandler(db *gorm.DB) UserHandler {
	return UserHandler{
		db: db,
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
	var sameEmail model.User
	err := uh.db.Where("email=?", ud.Email).First(&sameEmail).Error
	if err == nil {
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

// @Summary	Registers a user
// @Tags		users
// @Accept		json
// @Param		UserData	body	handler.RegisterReqData	true	"Data of user to register"
// @Produce	json
// @Success	201	{object}	handler.RegisterRespData
// @Failure	400	{object}	util.AppError
// @Failure	500	{object}	util.AppError
// @Router		/users/register [post]
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

	// err = service.SendMail()
	// if err != nil {
	// 	fmt.Println(err)
	// }

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

// @Summary	Login
// @Tags		users
// @Accept		json
// @Param		EmailPassword	body	handler.LoginReqData	true	"Email and password"
// @Produce	json
// @Success	200	{object}	handler.LoginRespData
// @Failure	404	{object}	util.AppError
// @Failure	400	{object}	util.AppError
// @Failure	500	{object}	util.AppError
// @Router		/users/login [post]
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
