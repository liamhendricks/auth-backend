package controllers

import (
	"fmt"

	"github.com/68696c6c/goat"
	"github.com/68696c6c/goat/query"
	"github.com/gin-gonic/gin"
	"github.com/liamhendricks/auth-backend/src/models"
	"github.com/liamhendricks/auth-backend/src/repos"
)

type DateController struct {
	dateRepo repos.DateRepo
	errors   goat.ErrorHandler
}

type datesResponse struct {
	Dates []*models.Date
}

type dateResponse struct {
	Date models.Date
}

type CreateDateRequest struct {
	Name string `json:"name" binding:"required"`
	Date string `json:"date" binding:"required"`
}

type UpdateDateRequest struct {
	Name string `json:"name"`
	Date string `json:"date"`
}

func NewDateController(dr repos.DateRepo, es goat.ErrorHandler) DateController {
	return DateController{
		dateRepo: dr,
		errors:   es,
	}
}

func (d *DateController) Store(c *gin.Context) {
	req, ok := goat.GetRequest(c).(*CreateDateRequest)
	if !ok {
		d.errors.HandleMessage(c, "failed to get request", goat.RespondBadRequestError)
		return
	}

	date := models.Date{
		Name: req.Name,
		Date: req.Date,
	}

	errs := d.dateRepo.Save(&date)
	if len(errs) > 0 {
		d.errors.HandleErrorsM(c, errs, "failed to save date record", goat.RespondBadRequestError)
		return
	}

	goat.RespondData(c, dateResponse{
		date,
	})
}

func (d *DateController) Index(c *gin.Context) {
	var q query.Query

	dates, errs := d.dateRepo.GetAll(&q)
	if len(errs) > 0 {
		goat.RespondServerErrors(c, errs)
		return
	}

	goat.RespondData(c, datesResponse{Dates: dates})
}

func (d *DateController) Show(c *gin.Context) {
	i := c.Param("id")
	id, err := goat.ParseID(i)
	if err != nil {
		d.errors.HandleErrorM(c, err, "failed to parse id: "+i, goat.RespondBadRequestError)
		return
	}

	date, errs := d.dateRepo.GetByID(id)
	if len(errs) > 0 {
		if goat.RecordNotFound(errs) {
			d.errors.HandleErrorsM(c, errs, "date record does not exist", goat.RespondNotFoundError)
			return
		} else {
			d.errors.HandleErrorsM(c, errs, "failed to get date record", goat.RespondServerError)
			return
		}
	}

	goat.RespondData(c, dateResponse{
		date,
	})
}

func (d *DateController) Update(c *gin.Context) {
	i := c.Param("id")
	id, err := goat.ParseID(i)
	if err != nil {
		d.errors.HandleErrorM(c, err, "failed to parse id: "+i, goat.RespondBadRequestError)
		return
	}

	req, ok := goat.GetRequest(c).(*UpdateDateRequest)
	if !ok {
		d.errors.HandleMessage(c, "failed to get request", goat.RespondBadRequestError)
		return
	}

	date, errs := d.dateRepo.GetByID(id)
	if len(errs) > 0 {
		if goat.RecordNotFound(errs) {
			d.errors.HandleErrorsM(c, errs, "date record does not exist", goat.RespondNotFoundError)
			return
		} else {
			d.errors.HandleErrorsM(c, errs, "failed to get date record", goat.RespondServerError)
			return
		}
	}

	if req.Name != "" {
		date.Name = req.Name
	}

	if req.Date != "" {
		date.Date = req.Date
	}

	errs = d.dateRepo.Save(&date)
	if len(errs) > 0 {
		d.errors.HandleErrorsM(c, errs, "failed to save date record", goat.RespondBadRequestError)
		return
	}

	goat.RespondMessage(c, fmt.Sprintf("%s has been updated", date.Name))
}

func (d *DateController) Delete(c *gin.Context) {
	i := c.Param("id")
	id, err := goat.ParseID(i)
	if err != nil {
		d.errors.HandleErrorM(c, err, "failed to parse id: "+i, goat.RespondBadRequestError)
		return
	}

	errs := d.dateRepo.Delete(id)
	if len(errs) > 0 {
		d.errors.HandleErrorsM(c, errs, "failed to delete date", goat.RespondBadRequestError)
		return
	}

	goat.RespondMessage(c, fmt.Sprintf("ID: %s has been deleted", id.String()))
}
