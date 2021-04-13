package services

import (
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
	APIKey             string
	BaseUrl            string
	ResetID            string
	PurchaseID         string
	SignupID           string
}

type TransactionalEmailType string

const Reset TransactionalEmailType = "Password Reset"
const Signup TransactionalEmailType = "Welcome"
const Purchase TransactionalEmailType = "Thanks for your purchase"

type TransactionalEmail struct {
	EmailType      TransactionalEmailType
	Data           map[string]string
	RecipientEmail string
	RecipientName  string
	ID             string
}

func NewSendgridMailer(email, name, baseUrl, key, resetID, purchaseID, signupID string) SendgridMailer {
	return SendgridMailer{
		SenderEmailAddress: email,
		SenderEmailName:    name,
		APIKey:             key,
		BaseUrl:            baseUrl,
		ResetID:            resetID,
		PurchaseID:         purchaseID,
		SignupID:           signupID,
	}
}

//send transactional email
func (m SendgridMailer) Send(emailObj *TransactionalEmail) error {
	emailSender := mail.NewV3Mail()
	e := mail.NewEmail(m.SenderEmailName, m.SenderEmailAddress)
	to := mail.NewEmail(emailObj.RecipientName, emailObj.RecipientEmail)
	emailSender.SetFrom(e)
	emailSender.SetTemplateID(emailObj.ID)
	p := mail.NewPersonalization()
	p.AddTos(to)

	for k, v := range emailObj.Data {
		p.SetDynamicTemplateData(k, v)
	}

	emailSender.AddPersonalizations(p)
	request := sendgrid.GetRequest(m.APIKey, "/v3/mail/send", "https://api.sendgrid.com")
	request.Method = "POST"
	var Body = mail.GetRequestBody(emailSender)
	request.Body = Body
	response, err := sendgrid.API(request)
	if err != nil || response.StatusCode > 301 {
		return err
	} else {
		return nil
	}
}

func (m SendgridMailer) CreateEmailOfType(data map[string]string, emailType TransactionalEmailType) *TransactionalEmail {
	var emailObj TransactionalEmail
	emailObj.Data = make(map[string]string)

	switch emailType {
	case Reset:
		reset(&emailObj, m.BaseUrl, m.ResetID, data)
		break
	case Purchase:
		purchase(&emailObj, m.BaseUrl, m.PurchaseID, data)
		break
	case Signup:
		signup(&emailObj, m.BaseUrl, m.SignupID, data)
		break
	}

	return &emailObj
}

func purchase(emailObj *TransactionalEmail, baseUrl, id string, data map[string]string) {
	emailObj.Data = data
	emailObj.ID = id
	emailObj.EmailType = Reset
	emailObj.RecipientEmail = data["email"]
	emailObj.RecipientName = data["name"]
}
func signup(emailObj *TransactionalEmail, baseUrl, id string, data map[string]string) {
	coursesLink := baseUrl + "/courses"
	galleryLink := baseUrl + "/gallery/all"
	data["coursesLink"] = coursesLink
	data["galleryLink"] = galleryLink
	emailObj.ID = id
	emailObj.Data = data
	emailObj.EmailType = Signup
	emailObj.RecipientEmail = data["email"]
	emailObj.RecipientName = data["name"]
}
func reset(emailObj *TransactionalEmail, baseUrl, id string, data map[string]string) {
	resetLink := baseUrl + "/app/reset?token=" + data["token"]
	data["resetLink"] = resetLink
	emailObj.ID = id
	emailObj.Data = data
	emailObj.EmailType = Reset
	emailObj.RecipientEmail = data["email"]
	emailObj.RecipientName = data["name"]
}
