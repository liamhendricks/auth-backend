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

func RequireAuth(e goat.ErrorHandler, ur repos.UserRepo, ss services.SessionService, pw services.PasswordService, requiredUserType models.UserType) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := strings.Split(c.Request.Header["Authorization"][0], " ")
		token := authHeader[0]
		id := authHeader[1]

		i, err := goat.ParseID(id)
		if err != nil {
			e.HandleContext(c, "access denied: failed to parse id "+i.String(), goat.RespondUnauthorizedError)
			return
		}

		user, errs := ur.GetByID(i, true)
		if len(errs) > 0 {
			e.HandleContext(c, "access denied: user not found", goat.RespondUnauthorizedError)
			return
		}

		t, err := goat.ParseID(token)
		if err != nil {
			e.HandleContext(c, "access denied: failed to parse token "+t.String(), goat.RespondUnauthorizedError)
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
