package common

import (
	"bytes"
	"ecom_be/configs"
	"fmt"
	"net/smtp"
	"text/template"
)

var (
	SMTP_HOST  = configs.GetEnv("SMTP_HOST")
	SMTP_PORT  = configs.GetEnv("SMTP_PORT")
	SMTP_EMAIL = configs.GetEnv("SMTP_EMAIL")
	SMTP_PASS  = configs.GetEnv("SMTP_PASSWORD")
	SMTP_NAME  = configs.GetEnv("SMTP_NAME")
)

var (
	ConfirmTemplate = "templates/confirm_account.html"
	// ... other templates
)

func TemplateHTML(templateHtml, data interface{}) (string, error) {

	dataFmt := struct {
		Name        string
		Title       string
		Description string
		Email       string
		Link        string
	}{
		Name:        data.(map[string]string)["name"],
		Title:       data.(map[string]string)["title"],
		Description: data.(map[string]string)["description"],
		Email:       data.(map[string]string)["email"],
		Link:        data.(map[string]string)["link"],
	}
	t, err := template.ParseFiles(templateHtml.(string))
	if err != nil {
		return "", err
	}

	var body bytes.Buffer
	if err := t.Execute(&body, dataFmt); err != nil {
		return "", err
	}

	return body.String(), nil
}

// SendMail sends an email using net/smtp
func SendMail(to []string, subject, body string) error {
	auth := smtp.PlainAuth("", SMTP_EMAIL, SMTP_PASS, SMTP_HOST)
	body, err := TemplateHTML(ConfirmTemplate, map[string]string{
		"name":        "Vu Cuong",
		"title":       "Confirm your account",
		"description": "Please click the link below to confirm your account:",
		"email":       to[0],
		"link":        "http://localhost:8080/api/v1/auth/confirm?token=your_token_here",
	})
	if err != nil {
		return err
	}

	msg := []byte(fmt.Sprintf(
		"From: %s\r\nTo: %s\r\nSubject: %s\r\nMIME-Version: 1.0\r\nContent-Type: text/html; charset=\"UTF-8\"\r\n\r\n%s",
		SMTP_NAME+" <"+SMTP_EMAIL+">",
		to[0], // Only first recipient in header
		subject,
		body,
	))
	addr := fmt.Sprintf("%s:%s", SMTP_HOST, SMTP_PORT)
	return smtp.SendMail(addr, auth, SMTP_EMAIL, to, msg)
}
