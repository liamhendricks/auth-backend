package middleware

import (
	"fmt"
	"strings"

	"github.com/68696c6c/goat"
	"github.com/gin-gonic/gin"
	"github.com/liamhendricks/auth-backend/src/models"
	"github.com/liamhendricks/auth-backend/src/repos"
	"github.com/liamhendricks/auth-backend/src/services"
)

func RequireAuth(e goat.ErrorHandler, ur repos.UserRepo, ss services.SessionService, requiredUserType models.UserType) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := strings.Split(c.Request.Header["Authorization"][0], " ")
		token := authHeader[0]
		id := authHeader[1]

		i, err := goat.ParseID(id)
		if err != nil {
			e.HandleContext(c, "access denied: failed to parse id "+id, goat.RespondUnauthorizedError)
			return
		}

		user, errs := ur.GetByID(i, true)
		if len(errs) > 0 {
			e.HandleContext(c, "access denied: user not found", goat.RespondUnauthorizedError)
			return
		}

		t, err := goat.ParseID(token)
		if err != nil {
			e.HandleContext(c, "access denied: failed to parse token "+token, goat.RespondUnauthorizedError)
			return
		}

		if !ss.Valid(&user, t) {
			e.HandleContext(c,
				fmt.Sprintf("access denied: token invalid for %s", user.Name),
				goat.RespondUnauthorizedError)
			return
		}

		if !user.UserType.IsGreaterThanEqTo(requiredUserType) {
			e.HandleContext(c, "access denied: not correct access", goat.RespondUnauthorizedError)
			return
		}
	}
}

func RequireSelf(e goat.ErrorHandler, ur repos.UserRepo, ss services.SessionService) gin.HandlerFunc {
	return func(c *gin.Context) {
		i := c.Param("id")

		pid, err := goat.ParseID(i)
		if err != nil {
			e.HandleContext(c, "access denied: failed to parse id "+i, goat.RespondUnauthorizedError)
			return
		}

		user, errs := ur.GetByID(pid, true)
		if len(errs) > 0 {
			e.HandleContext(c, "access denied: user not found", goat.RespondUnauthorizedError)
			return
		}

		authHeader := strings.Split(c.Request.Header["Authorization"][0], " ")
		token := authHeader[0]
		i = authHeader[1]

		hid, err := goat.ParseID(i)
		if err != nil {
			e.HandleContext(c, "access denied: failed to parse id "+i, goat.RespondUnauthorizedError)
			return
		}

		if hid != user.ID {
			e.HandleContext(c, "access denied: wrong user", goat.RespondUnauthorizedError)
			return
		}

		ht, err := goat.ParseID(token)
		if err != nil {
			e.HandleContext(c, "access denied: failed to parse id "+token, goat.RespondUnauthorizedError)
			return
		}

		if !ss.Valid(&user, ht) {
			e.HandleContext(c,
				fmt.Sprintf("access denied: token invalid for %s", user.Name),
				goat.RespondUnauthorizedError)
			return
		}
	}
}
