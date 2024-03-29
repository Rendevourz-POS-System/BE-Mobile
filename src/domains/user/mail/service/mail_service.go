package service

import (
	"fmt"
	"gopkg.in/gomail.v2"
	"log"
	"main.go/configs/app"
	"main.go/domains/user/mail"
	"os"
)

var sender mail.Service

type (
	// GmailSender Config for mail sender (Sender Payload Data) And Entity
	GmailSender struct {
		Host              string
		Port              int
		Name              string
		FromEmailAddress  string
		FromEmailPassword string
	}
)

func NewEmailSend() mail.Service {
	if sender == nil {
		sender = &GmailSender{
			Host:              app.GetConfig().Email.SenderHost,
			Port:              app.GetConfig().Email.SenderPort,
			Name:              app.GetConfig().Email.SenderEmailName,
			FromEmailAddress:  app.GetConfig().Email.SenderEmailAddress,
			FromEmailPassword: app.GetConfig().Email.SenderEmailPassword,
		}
	}
	return sender
}

func (service *GmailSender) ExecuteSendEmail(subject, content string, to, cc, bcc []string, attach string) error {
	dialer := gomail.NewDialer(
		service.Host,
		service.Port,
		service.FromEmailAddress,
		service.FromEmailPassword,
	)
	mailer := gomail.NewMessage()
	mailer.SetHeader("From", service.FromEmailAddress)
	mailer.SetHeader("To", to...)
	mailer.SetHeader("Cc", cc...)
	mailer.SetHeader("Bcc", bcc...)
	mailer.SetHeader("Subject", subject)
	mailer.SetBody("text/html",
		"Hello, <b>have a nice day </b>\n"+
			"<a href='"+content+"'> Here ! </a>\n"+
			"<h1>"+content+"</h1>"+
			"\n\n\n\nThank You For Using Shelter-Apps"+
			"Regards,\n Admin")
	if attach != "" {
		mailer.Attach(attach)
	}
	err := dialer.DialAndSend(mailer)
	if err != nil {
		a, _ := os.Getwd()
		fmt.Println("Dir Attach : ", attach)
		fmt.Println("DIR Config : ", app.GetConfig().Email.File)
		fmt.Println("Dir : ", a)
		log.Fatal(err.Error())
	}
	return nil
}