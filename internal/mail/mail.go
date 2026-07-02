package mail

import (
	"fmt"
	"net/smtp"
	"os"

	"github.com/joho/godotenv"
)

func init() {
	if err := godotenv.Load(); err != nil {
		fmt.Println("Warning: .env file tidak ditemukan")
	}
}

func sendEmail(toEmail string, subject string, htmlBody string) error {
	host := os.Getenv("SMTP_HOST") // smtp.gmail.com
	port := os.Getenv("SMTP_PORT") // 587
	user := os.Getenv("SMTP_USER")
	pass := os.Getenv("SMTP_PASS") // App Password Anda

	// Header wajib untuk Gmail agar tidak masuk Spam
	msg := []byte("From: " + user + "\r\n" +
		"To: " + toEmail + "\r\n" +
		"Subject: " + subject + "\r\n" +
		"MIME-version: 1.0;\r\n" +
		"Content-Type: text/html; charset=\"UTF-8\";\r\n" +
		"\r\n" +
		htmlBody)

	auth := smtp.PlainAuth("", user, pass, host)

	// Kirim email
	err := smtp.SendMail(host+":"+port, auth, user, []string{toEmail}, msg)
	
	if err != nil {
		fmt.Printf("Gmail SMTP Error: %v\n", err)
	}
	
	return err
}

func SendVerificationEmail(email string, token string) error {
	link := fmt.Sprintf("http://localhost:8080/verify-email?token=%s", token)
	body := fmt.Sprintf(`<h2>Verifikasi Email</h2><p>Klik tombol berikut untuk memverifikasi akun Anda:</p><a href="%s">Verify Email</a>`, link)
	
	return sendEmail(email, "Verify Your Email", body)
}

func SendResetPasswordEmail(email string, token string) error {
	link := fmt.Sprintf("http://localhost:8080/reset-password?token=%s", token)
	body := fmt.Sprintf(`<h2>Reset Password</h2><p>Klik tautan berikut untuk mereset password Anda:</p><a href="%s">Reset Password</a>`, link)
	
	return sendEmail(email, "Permintaan Reset Password", body)
}