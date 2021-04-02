package services

import (
	"fmt"

	"github.com/sendgrid/sendgrid-go"
	"github.com/sendgrid/sendgrid-go/helpers/mail"
)

type MailService interface {
	Send(emailObj *TransactionalEmail) error
	CreateEmailOfType(data map[string]string, emailType TransactionalEmailType) *TransactionalEmail
}

type SendgridMailer struct {
	SenderEmailAddress string
	SenderEmailName    string
	BaseUrl            string
	client             *sendgrid.Client
}

type TransactionalEmailType string

const Reset TransactionalEmailType = "Password Reset"
const Signup TransactionalEmailType = "Welcome"
const Purchase TransactionalEmailType = "Thanks for your purchase"

type TransactionalEmail struct {
	EmailType      TransactionalEmailType
	Data           string
	RecipientEmail string
	RecipientName  string
}

func NewSendgridMailer(email, name, baseUrl, key string) SendgridMailer {
	client := sendgrid.NewSendClient(key)
	return SendgridMailer{
		SenderEmailAddress: email,
		SenderEmailName:    name,
		BaseUrl:            baseUrl,
		client:             client,
	}
}

func (m SendgridMailer) Send(emailObj *TransactionalEmail) error {
	to := mail.NewEmail(emailObj.RecipientName, emailObj.RecipientEmail)
	from := mail.NewEmail(m.SenderEmailName, m.SenderEmailAddress)

	message := mail.NewSingleEmail(from, string(emailObj.EmailType), to, "", emailObj.Data)
	client := sendgrid.NewSendClient("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")
	response, err := client.Send(message)
	if err != nil || response.StatusCode > 301 {
		return err
	}

	return nil
}

func (m SendgridMailer) CreateEmailOfType(data map[string]string, emailType TransactionalEmailType) *TransactionalEmail {
	var emailObj TransactionalEmail

	switch emailType {
	case Reset:
		resetLink := m.BaseUrl + "/app/reset?token=" + data["token"]
		emailObj.Data = fmt.Sprintf("<p>Here is your password reset link for your art lessons account: <a href=%s>Link</a></p>", resetLink)
		emailObj.EmailType = Reset
		emailObj.RecipientEmail = data["email"]
		emailObj.RecipientName = data["name"]
		break
	case Purchase:
		profileLink := m.BaseUrl + "/app/profile"
		signupTemplate := `
<div>
  <h1>Thank you for purchasing: %s</h1>
  <p>Here is a link to your profile to take the course! <a href=%s>Profile</a></p>
  <p>Cheers, Lana Gloschat</p>
</div>
    `
		emailObj.Data = fmt.Sprintf(signupTemplate, data["course"], profileLink)
		emailObj.EmailType = Reset
		emailObj.RecipientEmail = data["email"]
		emailObj.RecipientName = data["name"]
		break
	case Signup:
		coursesLink := m.BaseUrl + "/courses"
		signupTemplate := `
<div>
  <h1>Thank you for signing up on my website!</h1>
  <p>I look forward to teaching you! Please check out all my courses <a href=%s>here</a>.</p>
  <p>Cheers, Lana Gloschat</p>
</div>
    `
		emailObj.Data = fmt.Sprintf(signupTemplate, coursesLink)
		emailObj.EmailType = Reset
		emailObj.RecipientEmail = data["email"]
		emailObj.RecipientName = data["name"]
		break
	}

	return &emailObj
}
