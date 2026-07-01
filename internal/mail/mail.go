package mail

import (
	"fmt"
	"os"
	"github.com/joho/godotenv"
	"gopkg.in/gomail.v2"
)

func init() {
	if err := godotenv.Load(); err != nil {
		fmt.Println("Warning: .env file tidak ditemukan, menggunakan environment variable sistem")
	}
}

// Host dan Port tetap const karena tidak berubah
const (
	SMTPHost = "smtp.gmail.com"
	SMTPPort = 587
)

// Helper untuk mengambil config
func getSMTPCreds() (string, string) {
	return os.Getenv("SMTP_USER"), os.Getenv("SMTP_PASS")
}

func SendVerificationEmail(email string, token string) error {

	smtpUser, smtpPass := getSMTPCreds()

	link := fmt.Sprintf(
		"https://huggingface.co/spaces/pasdaoiji/backend-oauth/verify-email?token=%s",
		token,
	)

	body := fmt.Sprintf(`
<h2>Verify Email</h2>

<p>Terima kasih sudah mendaftar.</p>

<p>Silakan klik tombol berikut:</p>

<a href="%s">Verify Email</a>

<p>Link berlaku 15 menit.</p>
`, link)

	m := gomail.NewMessage()

	m.SetHeader("From", smtpUser)
	m.SetHeader("To", email)
	m.SetHeader("Subject", "Verify Your Email")
	m.SetBody("text/html", body)

	d := gomail.NewDialer(
		SMTPHost,
		SMTPPort,
		smtpUser,
		smtpPass,
	)

	return d.DialAndSend(m)
}

func SendResetPasswordEmail(email string, token string) error {
	smtpUser, smtpPass := getSMTPCreds()
	// Sesuaikan URL dengan domain/endpoint reset password aplikasi Anda
	link := fmt.Sprintf(
		"https://huggingface.co/spaces/pasdaoiji/backend-oauth/tampilan/reset-password?token=%s",
		token,
	)

	body := fmt.Sprintf(`
<h2>Reset Password</h2>

<p>Kami menerima permintaan untuk mereset password akun Anda.</p>

<p>Silakan klik tautan di bawah ini untuk mengatur password baru:</p>

<a href="%s">Reset Password</a>

<p>Link ini hanya berlaku untuk 1 jam ke depan.</p>
<p>Jika Anda tidak merasa melakukan permintaan ini, silakan abaikan email ini.</p>
`, link)

	m := gomail.NewMessage()

	m.SetHeader("From", smtpUser)
	m.SetHeader("To", email)
	m.SetHeader("Subject", "Permintaan Reset Password")
	m.SetBody("text/html", body)

	d := gomail.NewDialer(
		SMTPHost,
		SMTPPort,
		smtpUser,
		smtpPass,
	)
	

	return d.DialAndSend(m)
}