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

// Fungsi helper untuk mengirim email via SMTP
func sendEmail(toEmail string, subject string, htmlBody string) error {
	host := os.Getenv("SMTP_HOST")    
	port := os.Getenv("SMTP_PORT")  
	user := os.Getenv("SMTP_USER")    
	pass := os.Getenv("SMTP_PASS")    
	// Header untuk memastikan email terformat sebagai HTML
	mime := "MIME-version: 1.0;\nContent-Type: text/html; charset=\"UTF-8\";\n\n"
	msg := []byte("Subject: " + subject + "\n" + mime + htmlBody)

	// Autentikasi ke server Brevo
	auth := smtp.PlainAuth("", user, pass, host)

	// Kirim email
	err := smtp.SendMail(host+":"+port, auth, user, []string{toEmail}, msg)
	return err
}

func SendVerificationEmail(email string, token string) error {
	link := fmt.Sprintf("https://oauth-go-backend-one.vercel.app/verify-email?token=%s", token)
	body := fmt.Sprintf(`<h2>Verify Email</h2><p>Klik tombol berikut:</p><a href="%s">Verify Email</a>`, link)
	
	return sendEmail(email, "Verify Your Email", body)
}

func SendResetPasswordEmail(email string, token string) error {
	link := fmt.Sprintf("https://oauth-go-backend-one.vercel.app/reset-password?token=%s", token)
	body := fmt.Sprintf(`<h2>Reset Password</h2><p>Klik tautan berikut:</p><a href="%s">Reset Password</a>`, link)
	
	return sendEmail(email, "Permintaan Reset Password", body)
}