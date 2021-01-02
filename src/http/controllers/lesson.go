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

type LessonController struct {
	lessonRepo repos.LessonRepo
	courseRepo repos.CourseRepo
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
	Name     string `json:"name" binding:"required"`
	CourseID string `json:"course_id" binding:"required"`
	Data     string `json:"data"`
}

type UpdateLessonRequest struct {
	Name     string `json:"name"`
	CourseID string `json:"course_id"`
	Data     string `json:"data"`
}

func NewLessonController(lr repos.LessonRepo, cr repos.CourseRepo, es goat.ErrorHandler) LessonController {
	return LessonController{
		lessonRepo: lr,
		courseRepo: cr,
		errors:     es,
	}
}

func (lc *LessonController) Index(c *gin.Context) {
	q := query.Query{}
	lessons, errs := lc.lessonRepo.GetAll(&q)
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
			lc.errors.HandleErrorsM(c, errs, "lesson does not exist", goat.RespondNotFoundError)
			return
		} else {
			lc.errors.HandleErrorsM(c, errs, "failed to get user", goat.RespondServerError)
			return
		}
	}

	goat.RespondData(c, lessonResponse{Lesson: lesson})
}

func (lc *LessonController) Store(c *gin.Context) {
	req, ok := goat.GetRequest(c).(*CreateLessonRequest)
	if !ok {
		lc.errors.HandleMessage(c, "failed to get request", goat.RespondBadRequestError)
		return
	}

	id, err := goat.ParseID(req.CourseID)
	if err != nil {
		lc.errors.HandleErrorM(c, err, "failed to parse ID", goat.RespondBadRequestError)
		return
	}

	lesson := models.Lesson{
		Name:     req.Name,
		CourseID: id,
		Data:     "{}",
	}

	if req.Data != "" {
		lesson.Data = req.Data
	}

	errs := lc.lessonRepo.Save(&lesson)
	if len(errs) > 0 {
		lc.errors.HandleErrorsM(c, errs, "failed to save lesson", goat.RespondBadRequestError)
		return
	}

	goat.RespondCreated(c, lessonResponse{Lesson: lesson})
}

func (lc *LessonController) Update(c *gin.Context) {
	i := c.Param("id")
	id, err := goat.ParseID(i)
	if err != nil {
		lc.errors.HandleErrorM(c, err, "failed to parse id: "+i, goat.RespondBadRequestError)
		return
	}

	req, ok := goat.GetRequest(c).(*UpdateLessonRequest)
	if !ok {
		lc.errors.HandleMessage(c, "failed to get request", goat.RespondBadRequestError)
		return
	}

	lesson, errs := lc.lessonRepo.GetByID(id)
	if len(errs) > 0 {
		if goat.RecordNotFound(errs) {
			lc.errors.HandleErrorsM(c, errs, "lesson does not exist", goat.RespondNotFoundError)
			return
		} else {
			lc.errors.HandleErrorsM(c, errs, "failed to get lesson", goat.RespondServerError)
			return
		}
	}

	lesson.Name = req.Name

	if req.Data != "" {
		lesson.Data = req.Data
	}

	courseID, err := goat.ParseID(req.CourseID)
	if err != nil {
		lc.errors.HandleErrorM(c, err, "failed to parse id: "+req.CourseID, goat.RespondBadRequestError)
		return
	}

	course, errs := lc.courseRepo.GetByID(courseID, false)
	if len(errs) > 0 {
		if goat.RecordNotFound(errs) {
			lc.errors.HandleErrorsM(c, errs, "course does not exist", goat.RespondNotFoundError)
			return
		} else {
			lc.errors.HandleErrorsM(c, errs, "failed to get course", goat.RespondServerError)
			return
		}
	}

	lesson.CourseID = course.ID

	errs = lc.lessonRepo.Save(&lesson)
	if len(errs) > 0 {
		lc.errors.HandleErrorsM(c, errs, "failed to save lesson", goat.RespondBadRequestError)
		return
	}

	goat.RespondMessage(c, fmt.Sprintf("ID: %s has been updated", lesson.ID.String()))
}

func (lc *LessonController) Delete(c *gin.Context) {
	i := c.Param("id")
	id, err := goat.ParseID(i)
	if err != nil {
		lc.errors.HandleErrorM(c, err, "failed to parse id: "+i, goat.RespondBadRequestError)
		return
	}

	errs := lc.lessonRepo.Delete(id)
	if len(errs) > 0 {
		lc.errors.HandleErrorsM(c, errs, "failed to delete user", goat.RespondBadRequestError)
	}

	goat.RespondMessage(c, fmt.Sprintf("ID: %s has been deleted", id.String()))
}
