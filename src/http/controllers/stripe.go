package controllers

import (
	"encoding/json"
	"errors"
	"io/ioutil"

	"github.com/68696c6c/goat"
	"github.com/gin-gonic/gin"
	"github.com/liamhendricks/auth-backend/src/models"
	"github.com/liamhendricks/auth-backend/src/repos"
	"github.com/liamhendricks/auth-backend/src/services"
	"github.com/stripe/stripe-go/v72"
	"github.com/stripe/stripe-go/v72/checkout/session"
	"github.com/stripe/stripe-go/v72/product"
	"github.com/stripe/stripe-go/v72/webhook"
)

type StripeController struct {
	errors  goat.ErrorHandler
	users   repos.UserRepo
	courses repos.CourseRepo
	mail    services.MailService
	secret  string
	key     string
}

func NewStripeController(
	es goat.ErrorHandler,
	secret string,
	key string,
	c repos.CourseRepo,
	m services.MailService,
	u repos.UserRepo) StripeController {
	return StripeController{
		errors:  es,
		secret:  secret,
		key:     key,
		users:   u,
		courses: c,
		mail:    m,
	}
}

func (s *StripeController) PaymentWebHook(c *gin.Context) {
	defer c.Request.Body.Close()

	body, err := ioutil.ReadAll(c.Request.Body)
	if err != nil {
		s.errors.HandleError(c, err, goat.RespondBadRequestError)
		return
	}

	event, err := webhook.ConstructEvent(body, c.GetHeader("Stripe-Signature"), s.secret)
	if err != nil {
		s.errors.HandleError(c, err, goat.RespondBadRequestError)
		return
	}

	switch event.Type {
	case models.CheckoutComplete.String():
		var checkoutSession stripe.CheckoutSession

		// get session object
		err = json.Unmarshal(event.Data.Raw, &checkoutSession)
		if err != nil {
			s.errors.HandleError(c, err, goat.RespondBadRequestError)
			return
		}

		params := &stripe.CheckoutSessionListLineItemsParams{
			Session: stripe.String(checkoutSession.ID),
		}

		stripe.Key = s.key

		// we need to request the line items from stripe
		i := session.ListLineItems(checkoutSession.ID, params)
		for i.Next() {
			li := i.LineItem()

			// we also need to retrieve the product so we can look at the metadata
			p, err := product.Get(li.Price.Product.ID, &stripe.ProductParams{})
			if err != nil {
				var errs []error
				e := errors.New("unable to get product from stripe")
				errs = append(errs, e, err)
				goat.RespondServerError(c, e)
				return
			}

			id, err := goat.ParseID(p.Metadata["course_server_uuid"])
			if err != nil {
				e := errors.New("unable to parse course_server_uuid from metadata")
				goat.RespondServerError(c, e)
				return
			}

			user, errs := s.users.GetByEmail(checkoutSession.CustomerEmail, true)
			if len(errs) > 0 {
				e := errors.New("unable to get user during webhook")
				errs = append(errs, e)
				goat.RespondServerErrors(c, errs)
				return
			}

			course, errs := s.courses.GetByID(id, false)
			if len(errs) > 0 {
				s.respondUser(&user, models.UserMissingCourse, c)
			}

			// now we can try to attach the course
			user.Courses = append(user.Courses, &course)
			user.UserType = models.PaidUser
			errs = s.users.Save(&user)
			if len(errs) > 0 {
				s.respondUser(&user, models.UserMissingCourse, c)
			}

			data := make(map[string]string)
			for _, c := range course.Dates {
				if c.Name == "critiqueDate" {
					data["critiqueDate"] = c.Date
				}

				if c.Name == "officeDate" {
					data["officeDate"] = c.Date
				}
			}
			data["email"] = user.Email
			data["name"] = user.Name
			data["course"] = course.Name
			data["facebookGroup"] = ""
			email := s.mail.CreateEmailOfType(data, services.Purchase)
			err = s.mail.Send(email)
			if err != nil {
				s.respondUser(&user, models.UserNoEmail, c)
			}
		}

		break
	}

	goat.RespondValid(c)
}

func (s *StripeController) respondUser(user *models.User, status models.UserStatus, c *gin.Context) {
	user.Status = status
	errs := s.users.Save(user)
	if len(errs) > 0 {
		goat.RespondServerErrors(c, errs)
		return
	}
	goat.RespondValid(c)
	return
}
