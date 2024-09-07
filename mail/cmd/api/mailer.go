package main

import (
	"bytes"
	"github.com/vanng822/go-premailer/premailer"
	mail "github.com/xhit/go-simple-mail/v2"
	"html/template"
	"time"
)

// Mail struct for mailer config
type Mail struct {
	Domain     string
	Host       string
	Port       int
	Username   string
	Password   string
	Encryption string
	SenderName string
	Sender     string
}

// Message struct for message
type Message struct {
	From        string
	FromName    string
	To          string
	Subject     string
	Attachments []string
	Data        any
	DataMap     map[string]any
}

// SendSMTPMessage will send an email using smtp
func (m *Mail) SendSMTPMessage(msg Message) error {
	// set sender
	if msg.From == "" {
		msg.From = m.Sender
	}

	// set sender name
	if msg.FromName == "" {
		msg.FromName = m.SenderName
	}

	// build html message
	data := map[string]any{
		"message": msg.Data,
	}

	msg.DataMap = data

	formattedMessage, err := m.buildHTMLMessage(msg)
	if err != nil {
		return err
	}

	plainMessage, err := m.buildPlainTextMessage(msg)
	if err != nil {
		return err
	}

	// setup smtp server
	server := mail.NewSMTPClient()
	server.Host = m.Host
	server.Port = m.Port
	server.Username = m.Username
	server.Password = m.Password
	server.Encryption = m.getEncryption(m.Encryption)
	server.KeepAlive = false
	server.ConnectTimeout = 10 * time.Second
	server.SendTimeout = 10 * time.Second

	smtpClient, err := server.Connect()
	if err != nil {
		return err
	}

	// send email
	email := mail.NewMSG()
	email.SetFrom(msg.From).AddTo(msg.To).SetSubject(msg.Subject)
	email.SetBody(mail.TextPlain, plainMessage)
	email.AddAlternative(mail.TextHTML, formattedMessage)

	if len(msg.Attachments) > 0 {
		for _, x := range msg.Attachments {
			email.AddAttachment(x)
		}
	}

	err = email.Send(smtpClient)
	if err != nil {
		return err
	}

	return nil
}

// getEncryption will get encryption type
func (m *Mail) getEncryption(s string) mail.Encryption {
	switch s {
	case "tls":
		return mail.EncryptionSTARTTLS
	case "ssl":
		return mail.EncryptionSSLTLS
	case "none", "":
		return mail.EncryptionNone
	default:
		return mail.EncryptionSTARTTLS
	}
}

// buildHTMLMessage will build plain text message
func (m *Mail) buildPlainTextMessage(msg Message) (string, error) {
	// templates location
	templateToRender := "./templates/mail.plain.gohtml"

	// get template for email
	t, err := template.New("email-plain").Parse(templateToRender)
	if err != nil {
		return "", err
	}

	var tpl bytes.Buffer
	if err = t.ExecuteTemplate(&tpl, "body", msg); err != nil {
		return "", err
	}

	return tpl.String(), nil
}

// buildHTMLMessage will build html message
func (m *Mail) buildHTMLMessage(msg Message) (string, error) {
	// templates location
	templateToRender := "./templates/mail.html.gohtml"

	// get template for email
	t, err := template.New("email-html").Parse(templateToRender)
	if err != nil {
		return "", err
	}

	var tpl bytes.Buffer
	if err = t.ExecuteTemplate(&tpl, "body", msg); err != nil {
		return "", err
	}

	formattedMessage := tpl.String()

	formattedMessage, err = m.inlineCSS(formattedMessage)
	if err != nil {
		return "", err
	}

	return formattedMessage, nil
}

// inlineCSS will inline css from html template
func (m *Mail) inlineCSS(s string) (string, error) {
	options := premailer.Options{
		RemoveClasses:     false,
		CssToAttributes:   false,
		KeepBangImportant: true,
	}

	prem, err := premailer.NewPremailerFromString(s, &options)
	if err != nil {
		return "", err
	}

	html, err := prem.Transform()
	if err != nil {
		return "", err
	}

	return html, nil
}
