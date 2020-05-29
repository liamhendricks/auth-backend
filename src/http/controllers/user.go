package controllers

import (
	"github.com/gin-gonic/gin"
	"github.com/liamhendricks/auth-backend/src/repos"
)

type UserController struct {
	userRepo repos.UserRepo
}

func NewUserController(repo repos.UserRepo) *UserController {
	return &UserController{
		userRepo: repo,
	}
}

func (u *UserController) Index(c *gin.Context) {
}

func (u *UserController) Show(c *gin.Context) {
}

func (u *UserController) Update(c *gin.Context) {
}

func (u *UserController) Delete(c *gin.Context) {
}

func (u *UserController) Store(c *gin.Context) {
}
