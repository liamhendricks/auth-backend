package routes

import (
	"github.com/68696c6c/goat"
	"github.com/68696c6c/goat/src/http"
	"github.com/gin-gonic/gin"
	"github.com/liamhendricks/auth-backend/src/app"
	"github.com/liamhendricks/auth-backend/src/http/controllers"
)

func InitRoutes(router http.Router, container app.ServiceContainer) {
	userController := controllers.NewUserController(container.UserRepo)
	engine := router.GetEngine()
	engine.GET("/health", Health)

	//users
	users := engine.Group("/users")
	users.GET("", userController.Index)
	users.POST("", userController.Store)
	users.GET("/:id", userController.Show)
	users.POST("/:id", userController.Update)
	users.DELETE("/:id", userController.Delete)
}

func Health(c *gin.Context) {
	goat.RespondMessage(c, "ok")
}
