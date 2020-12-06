package controllers

import (
	"fmt"
	"time"

	"github.com/68696c6c/goat"
	"github.com/68696c6c/goat/query"
	"github.com/gin-gonic/gin"
	"github.com/liamhendricks/auth-backend/src/models"
	"github.com/liamhendricks/auth-backend/src/repos"
	"github.com/liamhendricks/auth-backend/src/services"
)

type UserController struct {
	userRepo   repos.UserRepo
	resetRepo  repos.ResetRepo
	courseRepo repos.CourseRepo
	password   services.PasswordService
	errors     goat.ErrorHandler
}

func NewUserController(
	ur repos.UserRepo,
	rr repos.ResetRepo,
	cr repos.CourseRepo,
	ps services.PasswordService,
	es goat.ErrorHandler) UserController {
	return UserController{
		userRepo:   ur,
		courseRepo: cr,
		password:   ps,
		errors:     es,
	}
}

type CreateUserRequest struct {
	Name     string `json:"username" binding:"required"`
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
	UserType string `json:"user_type" binding:"required"`
}

type UpdateUserRequest struct {
	Name     string `json:"username" binding:"required"`
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
	UserType string `json:"user_type" binding:"required"`
}

type AttachCourseRequest struct {
	CourseName string `json:"course_name"`
}

type RevokeCourseRequest struct {
	CourseName string `json:"course_name"`
}

type ResetPasswordRequest struct {
	Token    string `json:"token" binding:"required"`
	Password string `json:"password" binding:"required"`
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
		p, err := u.password.Hash([]byte(req.Password))
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

	p, err := u.password.Hash([]byte(req.Password))
	if err != nil {
		u.errors.HandleErrorM(c, err, "cant handle password", goat.RespondServerError)
		return
	}

	userType, err := models.UserTypeFromString(req.UserType)
	if err != nil {
		u.errors.HandleErrorM(c, err, "incorrect user type", goat.RespondBadRequestError)
		return
	}

	user := models.User{
		Name:     req.Name,
		Email:    req.Email,
		UserType: userType,
		Password: string(p),
	}

	errs := u.userRepo.Save(&user)
	if len(errs) > 0 {
		u.errors.HandleErrorsM(c, errs, "failed to save user", goat.RespondServerError)
		return
	}

	goat.RespondCreated(c, userResponse{User: user})
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

	reset := models.Reset{
		Token:      goat.NewID(),
		UserID:     id,
		Expiration: time.Now().Add(30 * time.Minute),
	}

	errs = u.resetRepo.Save(&reset)
	if len(errs) > 0 {
		u.errors.HandleErrorsM(c, errs, "failed to save reset", goat.RespondServerError)
		return
	}

	//TODO: email service
	//fire off email with link

	goat.RespondMessage(c, fmt.Sprintf("reset password link has been sent to %s", user.Email))
}

func (u *UserController) ResetPassword(c *gin.Context) {
	i := c.Param("id")
	id, err := goat.ParseID(i)
	if err != nil {
		u.errors.HandleErrorM(c, err, "failed to parse id: "+i, goat.RespondBadRequestError)
		return
	}

	req, ok := goat.GetRequest(c).(*ResetPasswordRequest)
	if !ok {
		u.errors.HandleMessage(c, "failed to get request", goat.RespondBadRequestError)
		return
	}

	token, err := goat.ParseID(req.Token)
	if err != nil {
		u.errors.HandleErrorM(c, err, "failed to parse id: "+i, goat.RespondBadRequestError)
		return
	}

	reset, errs := u.resetRepo.GetByTokenUser(token, id)
	if len(errs) > 0 {
		if goat.RecordNotFound(errs) {
			u.errors.HandleErrorsM(c, errs, "token or user mismatch", goat.RespondNotFoundError)
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

func (u *UserController) AttachCourse(c *gin.Context) {
	i := c.Param("id")
	id, err := goat.ParseID(i)
	if err != nil {
		u.errors.HandleErrorM(c, err, "failed to parse id: "+i, goat.RespondBadRequestError)
		return
	}

	req, ok := goat.GetRequest(c).(*AttachCourseRequest)
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

	course, errs := u.courseRepo.GetByName(req.CourseName, false)
	if len(errs) > 0 {
		if goat.RecordNotFound(errs) {
			u.errors.HandleErrorsM(c, errs, "course does not exist", goat.RespondNotFoundError)
			return
		} else {
			u.errors.HandleErrorsM(c, errs, "failed to get course", goat.RespondServerError)
			return
		}
	}

	user.Courses = append(user.Courses, &course)

	errs = u.userRepo.Save(&user)
	if len(errs) > 0 {
		u.errors.HandleErrorsM(c, errs, "failed to save user", goat.RespondServerError)
		return
	}

	goat.RespondMessage(c, fmt.Sprintf("%s has been added to %s's account", course.Name, user.Name))
}

func (u *UserController) RevokeCourse(c *gin.Context) {
	i := c.Param("id")
	id, err := goat.ParseID(i)
	if err != nil {
		u.errors.HandleErrorM(c, err, "failed to parse id: "+i, goat.RespondBadRequestError)
		return
	}

	req, ok := goat.GetRequest(c).(*RevokeCourseRequest)
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

	course, errs := u.courseRepo.GetByName(req.CourseName, false)
	if len(errs) > 0 {
		if goat.RecordNotFound(errs) {
			u.errors.HandleErrorsM(c, errs, "course does not exist", goat.RespondNotFoundError)
			return
		} else {
			u.errors.HandleErrorsM(c, errs, "failed to get course", goat.RespondServerError)
			return
		}
	}

	for k, c := range user.Courses {
		if c.Name == req.CourseName {
			user.Courses = revoke(k, user.Courses)
			break
		}
	}

	errs = u.userRepo.Save(&user)
	if len(errs) > 0 {
		u.errors.HandleErrorsM(c, errs, "failed to save user", goat.RespondServerError)
		return
	}

	goat.RespondMessage(c, fmt.Sprintf("%s has been added to %s's account", course.Name, user.Name))
}

func revoke(key int, courses []*models.Course) []*models.Course {
	courses[len(courses)-1], courses[key] = courses[key], courses[len(courses)-1]
	return courses[:len(courses)-1]
}
