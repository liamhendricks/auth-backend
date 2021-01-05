package controllers

import (
	"fmt"
	"time"

	"github.com/68696c6c/goat"
	"github.com/gin-gonic/gin"
	"github.com/liamhendricks/auth-backend/src/models"
	"github.com/liamhendricks/auth-backend/src/repos"
	"github.com/liamhendricks/auth-backend/src/services"
)

type AuthController struct {
	userRepo        repos.UserRepo
	passWordService services.PasswordService
	sessionService  services.SessionService
	errors          goat.ErrorHandler
}

func NewAuthController(
	ur repos.UserRepo,
	ps services.PasswordService,
	ss services.SessionService,
	es goat.ErrorHandler) AuthController {
	return AuthController{
		userRepo:        ur,
		passWordService: ps,
		sessionService:  ss,
		errors:          es,
	}
}

type LoginRequest struct {
	Password string `json:"password" binding:"required"`
	Name     string `json:"username"`
	Email    string `json:"email"`
}

type LogoutRequest struct {
	Name         string `json:"username" binding:"required"`
	SessionToken string `json:"session_token" binding:"required"`
}

type CheckSessionRequest struct {
	Name         string `json:"username" binding:"required"`
	SessionToken string `json:"session_token" binding:"required"`
}

//TODO: session should be handled with HttpOnly cookie
type sessionResponse struct {
	models.User
}

func (ac *AuthController) Login(c *gin.Context) {
	req, ok := goat.GetRequest(c).(*LoginRequest)
	if !ok {
		ac.errors.HandleMessage(c, "failed to get request", goat.RespondBadRequestError)
		return
	}

	user, errs := ac.userRepo.GetByNameOrEmail(req.Name, req.Email, true)
	if len(errs) > 0 {
		if goat.RecordNotFound(errs) {
			ac.errors.HandleContext(c,
				fmt.Sprintf("could not find user with name: %s or email: %s", req.Name, req.Email),
				goat.RespondAuthenticationError)
			return
		} else {
			ac.errors.HandleMessage(c, "failed to get user", goat.RespondServerError)
			return
		}
	}

	valid := ac.passWordService.Compare(user.Password, req.Password)
	if valid {
		err := ac.sessionService.Start(&user)
		if err != nil {
			ac.errors.HandleMessage(c, "error starting session", goat.RespondServerError)
			return
		}

		goat.RespondData(c, sessionResponse{
			user,
		})
		return
	} else {
		ac.errors.HandleContext(c, "password incorrect", goat.RespondAuthenticationError)
		return
	}
}

func (ac *AuthController) Logout(c *gin.Context) {
	req, ok := goat.GetRequest(c).(*LogoutRequest)
	if !ok {
		ac.errors.HandleMessage(c, "failed to get request", goat.RespondBadRequestError)
		return
	}

	user, errs := ac.userRepo.GetByNameOrEmail(req.Name, "", true)
	if len(errs) > 0 {
		if goat.RecordNotFound(errs) {
			ac.errors.HandleContext(c,
				fmt.Sprintf("could not find user with name: %s", req.Name),
				goat.RespondAuthenticationError)
			return
		} else {
			ac.errors.HandleMessage(c, "failed to get user", goat.RespondServerError)
			return
		}
	}

	err := ac.sessionService.Stop(&user)
	if err != nil {
		ac.errors.HandleMessage(c, "error handling session", goat.RespondServerError)
		return
	}
}

func (ac *AuthController) CheckSession(c *gin.Context) {
	req, ok := goat.GetRequest(c).(*CheckSessionRequest)
	if !ok {
		ac.errors.HandleMessage(c, "failed to get request", goat.RespondBadRequestError)
		return
	}

	user, errs := ac.userRepo.GetByNameOrEmail(req.Name, "", true)
	if len(errs) > 0 {
		if goat.RecordNotFound(errs) {
			ac.errors.HandleContext(c,
				fmt.Sprintf("could not find user with name: %s", req.Name),
				goat.RespondAuthenticationError)
			return
		} else {
			ac.errors.HandleMessage(c, "failed to get user", goat.RespondServerError)
			return
		}
	}

	token, err := goat.ParseID(req.SessionToken)
	if err != nil {
		ac.errors.HandleErrorM(c, err, "failed to parse token "+token.String(), goat.RespondBadRequestError)
		return
	}

	if !ac.sessionService.Valid(&user, token) {
		ac.errors.HandleContext(c,
			fmt.Sprintf("could not validate session for %s", user.Name),
			goat.RespondAuthenticationError)
		return
	}

	session, err := ac.sessionService.Refresh(&user, 30*time.Minute)
	if err != nil {
		ac.errors.HandleErrorM(c, err, "error refreshing session", goat.RespondServerError)
		return
	}

	user.Session = session

	goat.RespondData(c, sessionResponse{
		user,
	})
}
