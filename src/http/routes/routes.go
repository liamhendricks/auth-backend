package routes

import (
	"github.com/68696c6c/goat"
	"github.com/68696c6c/goat/src/http"
	"github.com/gin-gonic/gin"
	"github.com/liamhendricks/auth-backend/src/app"
	"github.com/liamhendricks/auth-backend/src/http/controllers"
)

func InitRoutes(router http.Router, c app.ServiceContainer) {
	userController := controllers.NewUserController(c.UserRepo, c.PasswordService, c.Errors)
	lessonController := controllers.NewLessonController(c.LessonRepo, c.Errors)
	engine := router.GetEngine()
	engine.GET("/health", Health)

	//users
	users := engine.Group("/users")
	users.GET("", userController.Index)
	users.GET("/:id", userController.Show)
	users.POST("", goat.BindRequestMiddleware(controllers.CreateUserRequest{}), userController.Store)
	users.POST("/:id", goat.BindRequestMiddleware(controllers.UpdateUserRequest{}), userController.Update)
	users.DELETE("/:id", userController.Delete)
	users.GET("/:id/lessons/:id")
	users.GET("/:id/lessons/:id/revoke")

	lessons := engine.Group("/lessons")
	lessons.GET("", lessonController.Index)
	lessons.GET("/:id", lessonController.Show)
	lessons.POST("", goat.BindRequestMiddleware(controllers.CreateLessonRequest{}), lessonController.Store)
	lessons.POST("/:id", goat.BindRequestMiddleware(controllers.UpdateLessonRequest{}), lessonController.Update)
}

func Health(c *gin.Context) {
	goat.RespondMessage(c, "ok")
}
