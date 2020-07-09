package controllers

import (
	"github.com/68696c6c/goat"
	"github.com/gin-gonic/gin"
	"github.com/liamhendricks/auth-backend/src/repos"
	"github.com/liamhendricks/auth-backend/src/services"
)

type LessonController struct {
	lessonRepo repos.LessonRepo
	password   services.PasswordService
	errors     goat.ErrorHandler
}

func (lc *LessonController) Index(c *gin.Context)  {}
func (lc *LessonController) Show(c *gin.Context)   {}
func (lc *LessonController) Store(c *gin.Context)  {}
func (lc *LessonController) Update(c *gin.Context) {}
func (lc *LessonController) Delete(c *gin.Context) {}
