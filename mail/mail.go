package mail

import (
	"net/smtp"
	"strings"
)

const (
	HOST        = "smtp.163.com"
	SERVER_ADDR = "smtp.163.com:25"
	USER        = "16620808100@163.com" //发送邮件的邮箱
	PASSWORD    = "w123123"             //发送邮件邮箱的密码
)

type Email struct {
	to       string "to"
	subject  string "subject"
	msg      string "msg"
	mailtype string "html"
}

func NewEmail(to, subject, msg, mailtype string) *Email {
	return &Email{to: to, subject: subject, msg: msg, mailtype: mailtype}
}

func SendEmail(email *Email) error {
	auth := smtp.PlainAuth("", USER, PASSWORD, HOST)
	sendTo := strings.Split(email.to, ";")
	done := make(chan error, 1024)
	var content_type string
	if email.mailtype == "html" {
		content_type = "Content-Type: text/" + email.mailtype + "; charset=UTF-8"
	} else {
		content_type = "Content-Type: text/plain" + "; charset=UTF-8"
	}

	go func() {
		defer close(done)
		for _, v := range sendTo {

			str := strings.Replace("From: "+USER+"~To: "+v+"~Subject: "+email.subject+"~"+content_type+"~~", "~", "\r\n", -1) + email.msg

			err := smtp.SendMail(
				SERVER_ADDR,
				auth,
				USER,
				[]string{v},
				[]byte(str),
			)
			done <- err
		}
	}()

	for i := 0; i < len(sendTo); i++ {
		<-done
	}

	return nil
}
