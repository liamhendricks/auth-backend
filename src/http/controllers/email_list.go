package controllers

import (
	"github.com/68696c6c/goat"
	"github.com/68696c6c/goat/query"
	"github.com/gin-gonic/gin"
	"github.com/liamhendricks/auth-backend/src/models"
	"github.com/liamhendricks/auth-backend/src/repos"
)

type EmailListController struct {
	emailListRepo repos.EmailListRepo
	errors        goat.ErrorHandler
}

func NewEmailListController(el repos.EmailListRepo, es goat.ErrorHandler) EmailListController {
	return EmailListController{
		errors:        es,
		emailListRepo: el,
	}
}

type CreateEmailListRequest struct {
	Name  string `json:"username" binding:"required"`
	Email string `json:"email" binding:"required"`
}

type emailListResponse struct {
	List []*models.EmailList
}

func (elc *EmailListController) Index(c *gin.Context) {
	list, errs := elc.emailListRepo.GetAll(&query.Query{})
	if len(errs) > 0 {
		goat.RespondServerErrors(c, errs)
		return
	}

	goat.RespondData(c, emailListResponse{List: list})
}
