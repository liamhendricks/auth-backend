package controllers

import (
	"fmt"

	"github.com/68696c6c/goat"
	"github.com/68696c6c/goat/query"
	"github.com/gin-gonic/gin"
	"github.com/liamhendricks/auth-backend/src/models"
	"github.com/liamhendricks/auth-backend/src/repos"
	"github.com/liamhendricks/auth-backend/src/services"
)

type UserController struct {
	userRepo repos.UserRepo
	password services.PasswordService
	errors   goat.ErrorHandler
}

func NewUserController(ur repos.UserRepo, ps services.PasswordService, es goat.ErrorHandler) UserController {
	return UserController{
		userRepo: ur,
		password: ps,
		errors:   es,
	}
}

type CreateUserRequest struct {
	models.User
}

type UpdateUserRequest struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type usersResponse struct {
	Users []*models.User
}

type userResponse struct {
	User models.User
}

func (u *UserController) Index(c *gin.Context) {
	users, errs := u.userRepo.GetAll(&query.Query{})
	if len(errs) > 0 {
		goat.RespondServerErrors(c, errs)
		return
	}

	goat.RespondData(c, usersResponse{Users: users})
}

func (u *UserController) Show(c *gin.Context) {
	i := c.Param("id")
	id, err := goat.ParseID(i)
	if err != nil {
		u.errors.HandleErrorM(c, err, "failed to parse id: "+i, goat.RespondBadRequestError)
		return
	}

	user, errs := u.userRepo.GetByID(id, false)
	if len(errs) > 0 {
		if goat.RecordNotFound(errs) {
			u.errors.HandleErrorsM(c, errs, "user does not exist", goat.RespondNotFoundError)
			return
		} else {
			u.errors.HandleErrorsM(c, errs, "failed to get user", goat.RespondServerError)
			return
		}
	}

	goat.RespondData(c, userResponse{User: user})
}

func (u *UserController) Update(c *gin.Context) {
	i := c.Param("id")
	id, err := goat.ParseID(i)
	if err != nil {
		u.errors.HandleErrorM(c, err, "failed to parse id: "+i, goat.RespondBadRequestError)
		return
	}

	user, errs := u.userRepo.GetByID(id, false)
	if len(errs) > 0 {
		if goat.RecordNotFound(errs) {
			u.errors.HandleErrorsM(c, errs, "user does not exist", goat.RespondNotFoundError)
			return
		} else {
			u.errors.HandleErrorsM(c, errs, "failed to get user", goat.RespondServerError)
			return
		}
	}

	req, ok := goat.GetRequest(c).(*UpdateUserRequest)
	if !ok {
		u.errors.HandleMessage(c, "failed to get request", goat.RespondBadRequestError)
		return
	}

	if req.Name != "" {
		user.Name = req.Name
	}

	if req.Email != "" {
		user.Email = req.Email
	}

	if req.Password != "" {
		p, err := u.password.Encrypt([]byte(req.Password))
		if err != nil {
			u.errors.HandleErrorM(c, err, "cant handle password", goat.RespondServerError)
			return
		}

		user.Password = string(p)
	}

	errs = u.userRepo.Save(&user)
	if len(errs) > 0 {
		u.errors.HandleErrorsM(c, errs, "failed to save user", goat.RespondBadRequestError)
		return
	}

	goat.RespondMessage(c, fmt.Sprintf("%s has been updated", user.Name))
}

func (u *UserController) Delete(c *gin.Context) {
	i := c.Param("id")
	id, err := goat.ParseID(i)
	if err != nil {
		u.errors.HandleErrorM(c, err, "failed to parse id: "+i, goat.RespondBadRequestError)
		return
	}

	errs := u.userRepo.Delete(id)
	if len(errs) > 0 {
		u.errors.HandleErrorsM(c, errs, "failed to delete user", goat.RespondBadRequestError)
	}

	goat.RespondMessage(c, fmt.Sprintf("ID: %s has been deleted", id.String()))
}

func (u *UserController) Store(c *gin.Context) {
	req, ok := goat.GetRequest(c).(*CreateUserRequest)
	if !ok {
		u.errors.HandleMessage(c, "failed to get request", goat.RespondBadRequestError)
		return
	}

	p, err := u.password.Encrypt([]byte(req.Password))
	if err != nil {
		u.errors.HandleErrorM(c, err, "cant handle password", goat.RespondServerError)
		return
	}

	user := models.User{
		Name:     req.Name,
		Email:    req.Email,
		Password: string(p),
	}

	errs := u.userRepo.Save(&user)
	if len(errs) > 0 {
		u.errors.HandleErrorsM(c, errs, "failed to save user", goat.RespondBadRequestError)
		return
	}

	goat.RespondCreated(c, userResponse{User: user})
}

func (u *UserController) AttachLesson(c *gin.Context) {

}

func (u *UserController) RevokeLesson(c *gin.Context) {

}
