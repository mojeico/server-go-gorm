package email

import (
	"bytes"
	"fmt"
	"html/template"
	"log"
	"net/smtp"

	"github.com/trucktrace/pkg/logger"

	"github.com/spf13/viper"
)

type Request struct {
	to      []string
	subject string
	body    string
}

func NewRequest(to []string, subject string) *Request {
	return &Request{
		to:      to,
		subject: subject,
	}
}

func (r *Request) sendMail() bool {

	MIME := "MIME-version: 1.0;\nContent-Type: text/html; charset=\"UTF-8\";\n\n"
	body := "To: " + r.to[0] + "\r\nSubject: " + r.subject + "\r\n" + MIME + "\r\n" + r.body
	SMTP := fmt.Sprintf("%s:%s", viper.GetString("email.host"), viper.GetString("email.port"))
	auth := smtp.PlainAuth("", viper.GetString("email.user"), viper.GetString("email.password"), viper.GetString("email.host"))
	if err := smtp.SendMail(SMTP, auth, viper.GetString("email.user"), r.to, []byte(body)); err != nil {
		return false
	}
	return true
}

func (r *Request) parseTemplate(fileName string, data interface{}) error {
	t, err := template.ParseFiles(fileName)
	if err != nil {
		logger.ErrorLogger("parseTemplate", "Cant parse template").Error("Error - " + err.Error())
		return err
	}
	buffer := new(bytes.Buffer)
	if err = t.Execute(buffer, data); err != nil {
		return err
	}
	r.body = buffer.String()
	return nil
}

func (r *Request) Send(templateName string, items interface{}) {
	err := r.parseTemplate(templateName, items)
	if err != nil {
		log.Fatal(err)
	}

	if ok := r.sendMail(); ok {
		logger.InfoLogger("Send").Info(fmt.Sprintf("Email has been sent to %s\n", r.to))
		return
	}

	logger.ErrorLogger("Send", fmt.Sprintf("Email was not send to %s\n", r.to)).Error("Email was not send ")
}
