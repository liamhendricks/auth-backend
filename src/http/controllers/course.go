package controllers

import (
	"fmt"
	"strconv"

	"github.com/68696c6c/goat"
	"github.com/68696c6c/goat/query"
	"github.com/gin-gonic/gin"
	"github.com/liamhendricks/auth-backend/src/models"
	"github.com/liamhendricks/auth-backend/src/repos"
)

type CourseController struct {
	courseRepo repos.CourseRepo
	lessonRepo repos.LessonRepo
	errors     goat.ErrorHandler
}

type coursesResponse struct {
	Courses []*models.Course
}

type courseResponse struct {
	Course models.Course
}

type CreateCourseRequest struct {
	Name string `json:"name" binding:"required"`
	Max  string `json:"max" binding:"required"`
	Type string `json:"course_type" binding:"required"`
}

type UpdateCourseRequest struct {
	Name string `json:"name"`
	Max  string `json:"max"`
	Type string `json:"course_type"`
}

type AttachLessonRequest struct {
	LessonName string `json:"lesson_name"`
}

func NewCourseController(
	cr repos.CourseRepo,
	lr repos.LessonRepo,
	es goat.ErrorHandler) CourseController {
	return CourseController{
		courseRepo: cr,
		lessonRepo: lr,
		errors:     es,
	}
}

func (cc *CourseController) Index(c *gin.Context) {
	var q query.Query
	qp := c.Request.URL.Query()

	if _, ok := qp["with_users"]; ok {
		q.Preload = append(q.Preload, "Users")
	}

	if _, ok := qp["with_lessons"]; ok {
		q.Preload = append(q.Preload, "Lessons")
	}

	courses, errs := cc.courseRepo.GetAll(&q)
	if len(errs) > 0 {
		goat.RespondServerErrors(c, errs)
		return
	}

	goat.RespondData(c, coursesResponse{Courses: courses})
}

func (cc *CourseController) Show(c *gin.Context) {
	i := c.Param("id")
	id, err := goat.ParseID(i)
	if err != nil {
		cc.errors.HandleErrorM(c, err, "failed to parse id: "+i, goat.RespondBadRequestError)
		return
	}

	course, errs := cc.courseRepo.GetByID(id, true)
	if len(errs) > 0 {
		if goat.RecordNotFound(errs) {
			cc.errors.HandleErrorsM(c, errs, "course does not exist", goat.RespondNotFoundError)
			return
		} else {
			cc.errors.HandleErrorsM(c, errs, "failed to get user", goat.RespondServerError)
			return
		}
	}

	goat.RespondData(c, courseResponse{Course: course})
}

func (cc *CourseController) Store(c *gin.Context) {
	req, ok := goat.GetRequest(c).(*CreateCourseRequest)
	if !ok {
		cc.errors.HandleMessage(c, "failed to get request", goat.RespondBadRequestError)
		return
	}

	t, err := models.CourseTypeFromString(req.Type)
	if err != nil {
		cc.errors.HandleErrorM(c, err, "failed to convert", goat.RespondBadRequestError)
		return
	}

	max, err := strconv.Atoi(req.Max)
	if err != nil {
		cc.errors.HandleErrorM(c, err, "failed to parse max", goat.RespondBadRequestError)
		return
	}

	course := models.Course{
		Name:       req.Name,
		CourseType: t,
		Max:        max,
	}

	errs := cc.courseRepo.Save(&course)
	if len(errs) > 0 {
		cc.errors.HandleErrorsM(c, errs, "failed to save course", goat.RespondBadRequestError)
		return
	}

	goat.RespondCreated(c, courseResponse{Course: course})
}

func (cc *CourseController) Update(c *gin.Context) {
	i := c.Param("id")
	id, err := goat.ParseID(i)
	if err != nil {
		cc.errors.HandleErrorM(c, err, "failed to parse id: "+i, goat.RespondBadRequestError)
		return
	}

	req, ok := goat.GetRequest(c).(*UpdateCourseRequest)
	if !ok {
		cc.errors.HandleMessage(c, "failed to get request", goat.RespondBadRequestError)
		return
	}

	course, errs := cc.courseRepo.GetByID(id, false)
	if len(errs) > 0 {
		if goat.RecordNotFound(errs) {
			cc.errors.HandleErrorsM(c, errs, "course does not exist", goat.RespondNotFoundError)
			return
		} else {
			cc.errors.HandleErrorsM(c, errs, "failed to get user", goat.RespondServerError)
			return
		}
	}

	course.Name = req.Name
	max, err := strconv.Atoi(req.Max)
	if err != nil {
		cc.errors.HandleErrorM(c, err, "failed to parse max", goat.RespondBadRequestError)
		return
	}

	course.Max = max

	courseType, err := models.CourseTypeFromString(req.Type)
	if err != nil {
		cc.errors.HandleErrorM(c, err, "failed to convert course type", goat.RespondBadRequestError)
		return
	}

	course.CourseType = courseType

	errs = cc.courseRepo.Save(&course)
	if len(errs) > 0 {
		cc.errors.HandleErrorsM(c, errs, "failed to save course", goat.RespondBadRequestError)
		return
	}

	goat.RespondCreated(c, courseResponse{Course: course})
}

func (cc *CourseController) Delete(c *gin.Context) {
	i := c.Param("id")
	id, err := goat.ParseID(i)
	if err != nil {
		cc.errors.HandleErrorM(c, err, "failed to parse id: "+i, goat.RespondBadRequestError)
		return
	}

	course, errs := cc.courseRepo.GetByID(id, true)
	if len(errs) > 0 {
		if goat.RecordNotFound(errs) {
			cc.errors.HandleErrorsM(c, errs, "course does not exist", goat.RespondNotFoundError)
			return
		} else {
			cc.errors.HandleErrorsM(c, errs, "failed to get user", goat.RespondServerError)
			return
		}
	}

	//clear association table
	cc.courseRepo.Clear(&course, "Users")

	errs = cc.courseRepo.Delete(id)
	if len(errs) > 0 {
		cc.errors.HandleErrorsM(c, errs, "failed to delete course", goat.RespondBadRequestError)
	}

	goat.RespondMessage(c, fmt.Sprintf("ID: %s has been deleted", id.String()))
}

func (cc *CourseController) AttachLesson(c *gin.Context) {
	i := c.Param("id")
	id, err := goat.ParseID(i)
	if err != nil {
		cc.errors.HandleErrorM(c, err, "failed to parse id: "+i, goat.RespondBadRequestError)
		return
	}

	req, ok := goat.GetRequest(c).(*AttachLessonRequest)
	if !ok {
		cc.errors.HandleMessage(c, "failed to get request", goat.RespondBadRequestError)
		return
	}

	course, errs := cc.courseRepo.GetByID(id, true)
	if len(errs) > 0 {
		if goat.RecordNotFound(errs) {
			cc.errors.HandleErrorsM(c, errs, "course does not exist", goat.RespondNotFoundError)
			return
		} else {
			cc.errors.HandleErrorsM(c, errs, "failed to get user", goat.RespondServerError)
			return
		}
	}

	lesson, errs := cc.lessonRepo.GetByName(req.LessonName)
	if len(errs) > 0 {
		if goat.RecordNotFound(errs) {
			cc.errors.HandleErrorsM(c, errs, "lesson does not exist", goat.RespondNotFoundError)
			return
		} else {
			cc.errors.HandleErrorsM(c, errs, "failed to get user", goat.RespondServerError)
			return
		}
	}

	for _, l := range course.Lessons {
		if l.Name == lesson.Name {
			goat.RespondMessage(c, fmt.Sprintf("%s is already attached to this course", lesson.Name))
			return
		}
	}

	course.Lessons = append(course.Lessons, &lesson)
	errs = cc.courseRepo.Save(&course)
	if len(errs) > 0 {
		cc.errors.HandleErrorsM(c, errs, "failed to save course", goat.RespondBadRequestError)
		return
	}

	goat.RespondCreated(c, courseResponse{Course: course})
}
