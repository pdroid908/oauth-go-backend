package mail

import (
	"fmt"

	"gopkg.in/gomail.v2"
)

const (
	SMTPHost = "smtp.gmail.com"
	SMTPPort = 587

	SMTPUser = "p1998nr@gmail.com"
	SMTPPass = "uxgy zeox ucco ejpo"
)

func SendVerificationEmail(email string, token string) error {

	link := fmt.Sprintf(
		"http://localhost:8080/verify-email?token=%s",
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

	m.SetHeader("From", SMTPUser)
	m.SetHeader("To", email)
	m.SetHeader("Subject", "Verify Your Email")
	m.SetBody("text/html", body)

	d := gomail.NewDialer(
		SMTPHost,
		SMTPPort,
		SMTPUser,
		SMTPPass,
	)

	return d.DialAndSend(m)
}

func SendResetPasswordEmail(email string, token string) error {
	// Sesuaikan URL dengan domain/endpoint reset password aplikasi Anda
	link := fmt.Sprintf(
		"http://localhost:3000/tampilan/reset-password?token=%s",
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

	m.SetHeader("From", SMTPUser)
	m.SetHeader("To", email)
	m.SetHeader("Subject", "Permintaan Reset Password")
	m.SetBody("text/html", body)

	d := gomail.NewDialer(
		SMTPHost,
		SMTPPort,
		SMTPUser,
		SMTPPass,
	)
	

	return d.DialAndSend(m)
}