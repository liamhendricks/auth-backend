package routes

import (
	"github.com/68696c6c/goat"
	"github.com/68696c6c/goat/src/http"
	"github.com/gin-gonic/gin"
	"github.com/liamhendricks/auth-backend/src/app"
	"github.com/liamhendricks/auth-backend/src/http/controllers"
	"github.com/liamhendricks/auth-backend/src/http/middleware"
	"github.com/liamhendricks/auth-backend/src/models"
)

func InitRoutes(router http.Router, c app.ServiceContainer) {
	userController := controllers.NewUserController(c.UserRepo,
		c.ResetRepo,
		c.CourseRepo,
		c.PasswordService,
		c.Errors)

	authController := controllers.NewAuthController(c.UserRepo,
		c.PasswordService,
		c.SessionService,
		c.Errors)

	lessonController := controllers.NewLessonController(c.LessonRepo, c.CourseRepo, c.Errors)
	courseController := controllers.NewCourseController(c.CourseRepo, c.LessonRepo, c.Errors)

	engine := router.GetEngine()
	engine.GET("/health", Health)

	//user endpoints, must be self
	users := engine.Group("/users")
	users.Use(middleware.RequireSelf(c.Errors, c.UserRepo, c.SessionService))
	{
		users.GET("/:id", userController.Show)
		users.GET("/:id/forgot", userController.ForgotPassword)
		users.DELETE("/:id", userController.Delete)

		users.GET("/:id/courses", userController.UserCourses)

		users.POST("/:id",
			goat.BindRequestMiddleware(controllers.UpdateUserRequest{}),
			userController.Update)

		users.POST("/:id/reset",
			goat.BindRequestMiddleware(controllers.ResetPasswordRequest{}),
			userController.ResetPassword)
	}

	//open (anyone can create users)
	users.POST("",
		goat.BindRequestMiddleware(controllers.CreateUserRequest{}),
		userController.Store)

	//TODO: these will be webhooks for stripe. figure out a good way to protect them.
	users.POST("/:id/courses/attach",
		goat.BindRequestMiddleware(controllers.AttachCourseRequest{}),
		userController.AttachCourse)

	users.POST("/:id/courses/revoke",
		goat.BindRequestMiddleware(controllers.RevokeCourseRequest{}),
		userController.RevokeCourse)

	//auth
	auth := engine.Group("/auth")
	auth.POST("/login",
		goat.BindRequestMiddleware(controllers.LoginRequest{}),
		authController.Login)
	auth.POST("/logout",
		goat.BindRequestMiddleware(controllers.LogoutRequest{}),
		authController.Logout)
	auth.POST("/check",
		goat.BindRequestMiddleware(controllers.CheckSessionRequest{}),
		authController.CheckSession)

	//api endpoints, must be admin
	api := engine.Group("/api")
	api.Use(middleware.RequireAuth(c.Errors, c.UserRepo, c.SessionService, models.AdminUser))
	{
		//users
		users = api.Group("/users")
		users.GET("", userController.Index)
		users.POST("/:id",
			goat.BindRequestMiddleware(controllers.UpdateUserRequest{}),
			userController.Update)
		users.POST("",
			goat.BindRequestMiddleware(controllers.CreateUserAPIRequest{}),
			userController.StoreAPI)
		users.DELETE("/:id", userController.Delete)
		users.POST("/:id/level", goat.BindRequestMiddleware(controllers.UserLevelRequest{}), userController.UserLevel)

		//lessons
		lessons := api.Group("/lessons")
		lessons.GET("", lessonController.Index)
		lessons.GET("/:id", lessonController.Show)
		lessons.DELETE("/:id", lessonController.Delete)

		lessons.POST("",
			goat.BindRequestMiddleware(controllers.CreateLessonRequest{}),
			lessonController.Store)

		lessons.POST("/:id",
			goat.BindRequestMiddleware(controllers.UpdateLessonRequest{}),
			lessonController.Update)

		//courses
		courses := api.Group("/courses")
		courses.GET("", courseController.Index)
		courses.GET("/:id", courseController.Show)
		courses.DELETE("/:id", courseController.Delete)

		courses.POST("",
			goat.BindRequestMiddleware(controllers.CreateCourseRequest{}),
			courseController.Store)

		courses.POST("/:id",
			goat.BindRequestMiddleware(controllers.UpdateCourseRequest{}),
			courseController.Update)

		courses.POST("/:id/lessons/attach",
			goat.BindRequestMiddleware(controllers.AttachLessonRequest{}),
			courseController.AttachLesson)
	}
}

func Health(c *gin.Context) {
	goat.RespondMessage(c, "ok")
}
