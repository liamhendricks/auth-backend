package controllers

import (
	"fmt"
	"time"

	"github.com/68696c6c/goat"
	"github.com/68696c6c/goat/query"
	"github.com/68696c6c/goat/query/filter"
	"github.com/gin-gonic/gin"
	"github.com/liamhendricks/auth-backend/src/models"
	"github.com/liamhendricks/auth-backend/src/repos"
	"github.com/liamhendricks/auth-backend/src/services"
)

type UserController struct {
	userRepo       repos.UserRepo
	resetRepo      repos.ResetRepo
	courseRepo     repos.CourseRepo
	password       services.PasswordService
	sessionService services.SessionService
	mailService    services.MailService
	errors         goat.ErrorHandler
}

func NewUserController(
	ur repos.UserRepo,
	rr repos.ResetRepo,
	cr repos.CourseRepo,
	ps services.PasswordService,
	ss services.SessionService,
	ms services.MailService,
	es goat.ErrorHandler) UserController {
	return UserController{
		userRepo:       ur,
		courseRepo:     cr,
		resetRepo:      rr,
		password:       ps,
		sessionService: ss,
		mailService:    ms,
		errors:         es,
	}
}

type CreateUserRequest struct {
	Name     string `json:"username" binding:"required"`
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type CreateUserAPIRequest struct {
	Name    string   `json:"name" binding:"required"`
	Email   string   `json:"email" binding:"required"`
	Courses []string `json:"Courses"`
}

type UpdateUserRequest struct {
	Name     string   `json:"name"`
	Email    string   `json:"email"`
	Courses  []string `json:"Courses"`
	Password string   `json:"password"`
}

type AttachCourseRequest struct {
	CourseName string `json:"course_name"`
}

type CreateForgotPasswordRequest struct {
	Email string `json:"email" binding:"required"`
}

type RevokeCourseRequest struct {
	CourseName string `json:"course_name"`
}

type ResetPasswordRequest struct {
	Token    string `json:"token" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type UserLevelRequest struct {
	UserLevel string `json:"user_level" binding:"required"`
}

type usersResponse struct {
	Users []*models.User
}

type userResponse struct {
	models.User
}

func (u *UserController) Index(c *gin.Context) {
	var q query.Query
	qp := c.Request.URL.Query()

	if _, ok := qp["with_courses"]; ok {
		q.Preload = append(q.Preload, "Courses")
	}

	users, errs := u.userRepo.GetAll(&q)
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

	user, errs := u.userRepo.GetByID(id, true)
	if len(errs) > 0 {
		if goat.RecordNotFound(errs) {
			u.errors.HandleErrorsM(c, errs, "user does not exist", goat.RespondNotFoundError)
			return
		} else {
			u.errors.HandleErrorsM(c, errs, "failed to get user", goat.RespondServerError)
			return
		}
	}

	goat.RespondData(c, userResponse{
		user,
	})
}

func (u *UserController) Update(c *gin.Context) {
	i := c.Param("id")
	id, err := goat.ParseID(i)
	if err != nil {
		u.errors.HandleErrorM(c, err, "failed to parse id: "+i, goat.RespondBadRequestError)
		return
	}

	req, ok := goat.GetRequest(c).(*UpdateUserRequest)
	if !ok {
		u.errors.HandleMessage(c, "failed to get request", goat.RespondBadRequestError)
		return
	}

	user, errs := u.userRepo.GetByID(id, true)
	if len(errs) > 0 {
		if goat.RecordNotFound(errs) {
			u.errors.HandleErrorsM(c, errs, "user does not exist", goat.RespondNotFoundError)
			return
		} else {
			u.errors.HandleErrorsM(c, errs, "failed to get user", goat.RespondServerError)
			return
		}
	}

	if req.Name != "" {
		user.Name = req.Name
	}

	if req.Email != "" {
		user.Email = req.Email
	}

	if req.Password != "" {
		p, err := u.password.Hash([]byte(req.Password))
		if err != nil {
			u.errors.HandleErrorM(c, err, "cant handle password", goat.RespondServerError)
			return
		}

		user.Password = string(p)
	}

	var cids []goat.ID

	//set user_courses
	for _, v := range req.Courses {
		cid, err := goat.ParseID(v)
		if err != nil {
			u.errors.HandleErrorM(c, err, "failed to parse id: "+v, goat.RespondBadRequestError)
			return
		}

		cids = append(cids, cid)
	}

	q := query.Query{}
	q.Filter = filter.NewFilter()
	q.WhereIn("id", cids)
	cs, errs := u.courseRepo.GetAll(&q)
	if len(errs) > 0 {
		goat.RespondServerErrors(c, errs)
		return
	}

	//clear all relationships and then just add the updated ones back
	u.userRepo.Clear(&user, "Courses")
	user.Courses = cs

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

	p, err := u.password.Hash([]byte(req.Password))
	if err != nil {
		u.errors.HandleErrorM(c, err, "cant handle password", goat.RespondServerError)
		return
	}

	user := models.User{
		Name:     req.Name,
		Email:    req.Email,
		UserType: models.FreeUser,
		Password: string(p),
	}

	errs := u.userRepo.Save(&user)
	if len(errs) > 0 {
		u.errors.HandleErrorsM(c, errs, "failed to save user", goat.RespondServerError)
		return
	}

	err = u.sessionService.Start(&user)
	if err != nil {
		u.errors.HandleErrorsM(c, errs, "failed to start session", goat.RespondServerError)
		return
	}

	data := make(map[string]string)
	data["email"] = user.Email
	data["name"] = user.Name

	email := u.mailService.CreateEmailOfType(data, services.Signup)
	err = u.mailService.Send(email)
	if err != nil {
		u.errors.HandleErrorsM(c, errs, "failed to send welcome email", goat.RespondServerError)
		return
	}

	goat.RespondData(c, userResponse{
		user,
	})
}

func (u *UserController) StoreAPI(c *gin.Context) {
	req, ok := goat.GetRequest(c).(*CreateUserAPIRequest)
	if !ok {
		u.errors.HandleMessage(c, "failed to get request", goat.RespondBadRequestError)
		return
	}

	//TODO: no
	pw := "password123"

	p, err := u.password.Hash([]byte(pw))
	if err != nil {
		u.errors.HandleErrorM(c, err, "cant handle password", goat.RespondServerError)
		return
	}

	user := models.User{
		Name:     req.Name,
		Email:    req.Email,
		UserType: models.FreeUser,
		Password: string(p),
	}

	errs := u.userRepo.Save(&user)
	if len(errs) > 0 {
		u.errors.HandleErrorsM(c, errs, "failed to save user", goat.RespondServerError)
		return
	}

	//TODO: send welcome email

	goat.RespondCreated(c, userResponse{
		user,
	})
}

func (u *UserController) UserCourses(c *gin.Context) {
	i := c.Param("id")
	id, err := goat.ParseID(i)
	if err != nil {
		u.errors.HandleErrorM(c, err, "failed to parse id: "+i, goat.RespondBadRequestError)
		return
	}

	user, errs := u.userRepo.GetByID(id, true)
	if len(errs) > 0 {
		if goat.RecordNotFound(errs) {
			u.errors.HandleErrorsM(c, errs, "user does not exist", goat.RespondNotFoundError)
			return
		} else {
			u.errors.HandleErrorsM(c, errs, "failed to get user", goat.RespondServerError)
			return
		}
	}

	goat.RespondCreated(c, coursesResponse{Courses: user.Courses})
}

func (u *UserController) ForgotPassword(c *gin.Context) {
	req, ok := goat.GetRequest(c).(*CreateForgotPasswordRequest)
	if !ok {
		u.errors.HandleMessage(c, "failed to get request", goat.RespondBadRequestError)
		return
	}

	user, errs := u.userRepo.GetByEmail(req.Email, false)
	if len(errs) > 0 {
		if goat.RecordNotFound(errs) {
			u.errors.HandleErrorsM(c, errs, "user does not exist", goat.RespondNotFoundError)
			return
		} else {
			u.errors.HandleErrorsM(c, errs, "failed to get user", goat.RespondServerError)
			return
		}
	}

	reset := &models.Reset{
		Token:      goat.NewID(),
		UserID:     user.ID,
		Expiration: time.Now().Add(30 * time.Minute),
	}

	errs = u.resetRepo.Save(reset)
	if len(errs) > 0 {
		u.errors.HandleErrorsM(c, errs, "failed to save reset", goat.RespondServerError)
		return
	}

	emailData := make(map[string]string)
	emailData["token"] = reset.Token.String()
	emailData["email"] = user.Email
	emailData["name"] = user.Name

	email := u.mailService.CreateEmailOfType(emailData, services.Reset)
	err := u.mailService.Send(email)
	if err != nil {
		u.errors.HandleErrorsM(c, errs, "failed to send email", goat.RespondBadRequestError)
		return
	}

	goat.RespondMessage(c, fmt.Sprintf("reset password link has been sent to %s", user.Email))
}

func (u *UserController) ResetPassword(c *gin.Context) {
	req, ok := goat.GetRequest(c).(*ResetPasswordRequest)
	if !ok {
		u.errors.HandleMessage(c, "failed to get request", goat.RespondBadRequestError)
		return
	}

	token, err := goat.ParseID(req.Token)
	if err != nil {
		u.errors.HandleErrorM(c, err, "failed to parse id: "+req.Token, goat.RespondBadRequestError)
		return
	}

	reset, errs := u.resetRepo.GetByToken(token)
	if len(errs) > 0 {
		if goat.RecordNotFound(errs) {
			u.errors.HandleErrorsM(c, errs, "token mismatch", goat.RespondNotFoundError)
			return
		} else {
			u.errors.HandleErrorsM(c, errs, "failed to get reset", goat.RespondServerError)
			return
		}
	}

	if time.Now().After(reset.Expiration) {
		u.errors.HandleMessage(c, "reset has expired", goat.RespondBadRequestError)
		return
	}

	user, errs := u.userRepo.GetByID(reset.UserID, true)
	if len(errs) > 0 {
		if goat.RecordNotFound(errs) {
			u.errors.HandleErrorsM(c, errs, "user does not exist", goat.RespondNotFoundError)
			return
		} else {
			u.errors.HandleErrorsM(c, errs, "failed to get user", goat.RespondServerError)
			return
		}
	}

	p, err := u.password.Hash([]byte(req.Password))
	if err != nil {
		u.errors.HandleErrorM(c, err, "cant handle password", goat.RespondServerError)
		return
	}

	user.Password = string(p)

	errs = u.userRepo.Save(&user)
	if len(errs) > 0 {
		u.errors.HandleErrorsM(c, errs, "failed to save user", goat.RespondServerError)
		return
	}

	errs = u.resetRepo.Delete(reset.ID)
	if len(errs) > 0 {
		u.errors.HandleErrorsM(c, errs, "failed to delete reset record", goat.RespondServerError)
		return
	}

	goat.RespondMessage(c, fmt.Sprintf("%s's password has been updated", user.Name))
}
