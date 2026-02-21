package services

import (
	"bytes"
	"fmt"
	"html/template"
	"strconv"

	"github.com/rendyfutsuy/base-go/utils"
	"gopkg.in/gomail.v2"
)

type EmailService struct {
	authEmail    string
	smtpHost     string
	authPassword string
	smtpPort     int
	senderEmail  string
	resetURL     string
}

func NewEmailService() (*EmailService, error) {
	port, err := strconv.Atoi(utils.ConfigVars.String("email.smtp_port"))
	if err != nil {
		return nil, err
	}

	return &EmailService{
		authEmail:    utils.ConfigVars.String("email.smtp_auth_email"),
		smtpHost:     utils.ConfigVars.String("email.smtp_host"),
		authPassword: utils.ConfigVars.String("email.smtp_password"),
		smtpPort:     port,
		senderEmail:  utils.ConfigVars.String("email.smtp_sender_mail"),
		resetURL:     utils.ConfigVars.String("email.reset_password_url"),
	}, nil
}

func (s *EmailService) SendPasswordResetEmail(email, session string) error {
	var tpl bytes.Buffer

	pathTemplate := "public/template/reset-password.html"
	subject := "Reset Password Link"

	tmpl, err := template.ParseFiles(pathTemplate)
	if err != nil {
		return err
	}

	data := map[string]interface{}{
		"reset_link": s.resetURL + "?token=" + session,
	}

	if err = tmpl.Execute(&tpl, data); err != nil {
		return err
	}

	m := gomail.NewMessage()
	m.SetHeader("From", s.senderEmail)
	m.SetHeader("To", email)
	m.SetHeader("Subject", subject)
	m.SetBody("text/html", tpl.String())

	d := gomail.NewDialer(s.smtpHost, s.smtpPort, s.authEmail, s.authPassword)
	return d.DialAndSend(m)
}

func (s *EmailService) SendVerificationEmail(email, code string) error {
	subject := "Verification Code"
	body := fmt.Sprintf("<p>Your verification code is: <strong>%s</strong></p>", code)

	m := gomail.NewMessage()
	m.SetHeader("From", s.senderEmail)
	m.SetHeader("To", email)
	m.SetHeader("Subject", subject)
	m.SetBody("text/html", body)

	d := gomail.NewDialer(s.smtpHost, s.smtpPort, s.authEmail, s.authPassword)
	return d.DialAndSend(m)
}
