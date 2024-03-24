package controller

import (
	"main.go/configs/app"
	Email "main.go/domains/user/entities"
	"main.go/domains/user/mail/service"
	"os"
	"path/filepath"
)

func SendEmail(send *Email.MailSend) error {
	cc := []string{}
	bcc := []string{}
	if send.Cc != "" {
		cc = []string{send.Cc}
	}
	if send.Bcc != "" {
		bcc = []string{send.Bcc}
	}
	if send.Attach != "" {
		dirPath, _ := os.Getwd()
		send.Attach = filepath.Join(dirPath, app.GetConfig().Email.File, send.Attach)
	}
	sender := service.NewEmailSend()
	err := sender.ExecuteSendEmail(
		send.Subject,
		send.Content,
		[]string{send.To},
		cc,
		bcc,
		send.Attach)
	if err != nil {
		return err
	}
	return nil
}
