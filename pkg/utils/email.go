package utils

import (
	"net/smtp"
	"os"
)

func SendVerificationEmail(email, link string) error {
	smtpHost := os.Getenv("SMTP_HOST")
	smtpPort := os.Getenv("SMTP_PORT")
	from := os.Getenv("SMTP_USER")
	password := os.Getenv("SMTP_PASSWORD")
	subject := "Email Verification"
	body := "Click the following link to verify your email: " + link

	auth := smtp.PlainAuth("", from, password, smtpHost)
	msg := []byte("Subject: " + subject + "\r\n\r\n" + body)
	return smtp.SendMail(smtpHost+":"+smtpPort, auth, from, []string{email}, msg)
}

func SendResetPasswordEmail(email, link string) error {
	smtpHost := os.Getenv("SMTP_HOST")
	smtpPort := os.Getenv("SMTP_PORT")
	from := os.Getenv("SMTP_USER")
	password := os.Getenv("SMTP_PASSWORD")
	subject := "Solicitud de restablecimiento de contraseña"
	body := "Haz clic en el siguiente enlace para restablecer tu contraseña: " + link

	auth := smtp.PlainAuth("", from, password, smtpHost)
	msg := []byte("Subject: " + subject + "\r\n\r\n" + body)
	return smtp.SendMail(smtpHost+":"+smtpPort, auth, from, []string{email}, msg)
}
