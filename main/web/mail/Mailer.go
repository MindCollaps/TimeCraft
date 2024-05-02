package mail

import (
	"bytes"
	"embed"
	"github.com/wneessen/go-mail"
	"io"
	"log"
	"os"
	"strconv"
	"text/template"
)

var files embed.FS
var Mailer *mail.Client
var MailFrom string
var NameFrom string

func InitMailer() bool {
	port, err := strconv.Atoi(os.Getenv("MAIL_PORT"))
	if err != nil {
		log.Print("Invalid Mailer port number")
		return false
	}

	host := os.Getenv("MAIL_HOST")
	username := os.Getenv("MAIL_USERNAME")
	password := os.Getenv("MAIL_PASSWORD")

	c, err := mail.NewClient(host, mail.WithPort(port), mail.WithUsername(username), mail.WithPassword(password))
	if err != nil {
		log.Print("Failed to create mail client - " + err.Error())
		return false
	}

	Mailer = c
	MailFrom = os.Getenv("MAIL_FROM")
	if MailFrom == "" {
		log.Print("MAIL_FROM not set in .env file - defaulting to 'info@timecraft.com'")
		MailFrom = "info@timecraft.com"
	}

	NameFrom = os.Getenv("MAIL_NAME")
	if NameFrom == "" {
		log.Print("MAIL_NAME not set in .env file - defaulting to 'TimeCraft'")
		NameFrom = "TimeCraft"
	}
	return true
}

func addBodyTemplate(mail mail.Msg, templateFile string, data any) bool {
	//Some black magic from template.SetBodyHTMLTemplate
	t := template.Must(template.ParseFS(files, "main/mail/templates/"+templateFile+".gohtml"))
	if t == nil {
		log.Println("Failed to parse template templateFile")
		return false
	}

	buf := bytes.Buffer{}
	if err := t.Execute(&buf, data); err != nil {
		log.Print("Failed to execute template - " + err.Error())
		return false
	}

	w := func(w io.Writer) (int64, error) {
		nb, err := w.Write(buf.Bytes())
		return int64(nb), err
	}

	mail.SetBodyWriter("text/html", w)
	return true
}

func defaultMail() *mail.Msg {
	m := mail.NewMsg()
	err := m.FromFormat(MailFrom, NameFrom)
	if err != nil {
		log.Print("Failed to set sender of mail - " + err.Error())
		return nil
	}

	err = m.EnvelopeFrom(MailFrom)
	if err != nil {
		log.Print("Failed to set envelope sender of mail - " + err.Error())
		return nil
	}

	m.SetMessageID()
	m.SetDate()

	return m
}

func sendDefaultTemplateMail(recipientEmail string, recipientName string, templateFile string, subject string, data any) {
	m := defaultMail()
	if m == nil {
		return
	}

	m.Subject(subject)
	err := m.AddToFormat(recipientEmail, recipientName)

	if err != nil {
		log.Print("Failed to add recipient to mail - " + err.Error())
		return
	}

	if !addBodyTemplate(*m, templateFile, data) {
		return
	}

	err = Mailer.DialAndSend(m)
	if err != nil {
		log.Print("Failed to send mail - " + err.Error())
		return
	}
}

func SendPasswordResetMail(recipientEmail string, recipientName string, resetLink string) {
	type PasswordResetData struct {
		Name string
		Link string
	}

	data := PasswordResetData{
		Name: recipientName,
		Link: resetLink,
	}

	sendDefaultTemplateMail(recipientEmail, recipientName, "passwordReset", "Password Reset", data)
}

func SendWelcomeEmail(recipientEmail string, recipientName string) {
	type WelcomeData struct {
		Name string
	}

	data := WelcomeData{
		Name: recipientName,
	}

	sendDefaultTemplateMail(recipientEmail, recipientName, "welcome", "Welcome to TimeCraft", data)
}
