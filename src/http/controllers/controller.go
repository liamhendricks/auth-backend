package controllers

import (
	"github.com/gin-gonic/gin"
)

type Controller interface {
	Index(c *gin.Context)
	Show(c *gin.Context)
	Store(c *gin.Context)
	Update(c *gin.Context)
	Delete(c *gin.Context)
}
