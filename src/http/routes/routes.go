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
		c.SessionService,
		c.MailService,
		c.Errors)
	authController := controllers.NewAuthController(c.UserRepo,
		c.PasswordService,
		c.SessionService,
		c.Errors)
	stripeController := controllers.NewStripeController(
		c.Errors,
		c.Config.StripeEndpointSecret,
		c.Config.StripeSecretKey,
		c.CourseRepo,
		c.MailService,
		c.UserRepo)
	lessonController := controllers.NewLessonController(c.LessonRepo, c.CourseRepo, c.Errors)
	courseController := controllers.NewCourseController(c.CourseRepo, c.LessonRepo, c.Errors)
	emailListController := controllers.NewEmailListController(c.EmailListRepo, c.Errors)

	engine := router.GetEngine()
	engine.GET("/health", Health)
	engine.POST("/forgot",
		goat.BindRequestMiddleware(controllers.CreateForgotPasswordRequest{}),
		userController.ForgotPassword)
	engine.POST("/reset",
		goat.BindRequestMiddleware(controllers.ResetPasswordRequest{}),
		userController.ResetPassword)
	engine.POST("/email-list", goat.BindRequestMiddleware(controllers.CreateEmailListRequest{}), emailListController.Store)

	//user endpoints, must be self
	users := engine.Group("/users")
	users.Use(middleware.RequireSelf(c.Errors, c.UserRepo, c.SessionService))
	{
		users.GET("/:id", userController.Show)
		users.GET("/:id/courses", userController.UserCourses)
		users.DELETE("/:id", userController.Delete)
		users.POST("/:id",
			goat.BindRequestMiddleware(controllers.UpdateUserRequest{}),
			userController.Update)
	}

	//open (anyone can create users)
	users = engine.Group("/users")
	users.POST("",
		goat.BindRequestMiddleware(controllers.CreateUserRequest{}),
		userController.Store)

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
		//email list
		list := api.Group("/email-list")
		list.GET("", emailListController.Index)

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

	webhooks := engine.Group("/webhooks")
	webhooks.POST("/stripe", stripeController.PaymentWebHook)
}

func Health(c *gin.Context) {
	goat.RespondMessage(c, "ok")
}
