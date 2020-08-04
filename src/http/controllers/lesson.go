package controllers

import (
	"github.com/68696c6c/goat"
	"github.com/68696c6c/goat/query"
	"github.com/gin-gonic/gin"
	"github.com/liamhendricks/auth-backend/src/models"
	"github.com/liamhendricks/auth-backend/src/repos"
	"github.com/liamhendricks/auth-backend/src/services"
)

type LessonController struct {
	lessonRepo repos.LessonRepo
	password   services.PasswordService
	errors     goat.ErrorHandler
}

type lessonsResponse struct {
	Lessons []*models.Lesson
}

type lessonResponse struct {
	Lesson models.Lesson
}

type CreateLessonRequest struct {
	Name       string `json:"name" binding:"required"`
	LessonType string `json:"lesson_type" binding:"required"`
}

type UpdateLessonRequest struct {
	Name       string `json:"name" binding:"required"`
	LessonType string `json:"lesson_type" binding:"required"`
}

func NewLessonController(lr repos.LessonRepo, es goat.ErrorHandler) LessonController {
	return LessonController{
		lessonRepo: lr,
		errors:     es,
	}
}

func (lc *LessonController) Index(c *gin.Context) {
	lessons, errs := lc.lessonRepo.GetAll(&query.Query{})
	if len(errs) > 0 {
		goat.RespondServerErrors(c, errs)
		return
	}

	goat.RespondData(c, lessonsResponse{Lessons: lessons})
}

func (lc *LessonController) Show(c *gin.Context) {
	i := c.Param("id")
	id, err := goat.ParseID(i)
	if err != nil {
		lc.errors.HandleErrorM(c, err, "failed to parse id: "+i, goat.RespondBadRequestError)
		return
	}

	lesson, errs := lc.lessonRepo.GetByID(id)
	if len(errs) > 0 {
		if goat.RecordNotFound(errs) {
			lc.errors.HandleErrorsM(c, errs, "user does not exist", goat.RespondNotFoundError)
			return
		} else {
			lc.errors.HandleErrorsM(c, errs, "failed to get user", goat.RespondServerError)
			return
		}
	}

	goat.RespondData(c, lessonResponse{Lesson: lesson})
}

func (lc *LessonController) Store(c *gin.Context) {
	_, ok := goat.GetRequest(c).(*CreateLessonRequest)
	if !ok {
		lc.errors.HandleMessage(c, "failed to get request", goat.RespondBadRequestError)
		return
	}
}
func (lc *LessonController) Update(c *gin.Context) {}
func (lc *LessonController) Delete(c *gin.Context) {}
