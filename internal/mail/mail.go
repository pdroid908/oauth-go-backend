package mail

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
	"github.com/resend/resend-go/v2"
)

func init() {
	if err := godotenv.Load(); err != nil {
		fmt.Println("Warning: .env file tidak ditemukan")
	}
}

// Fungsi helper untuk inisialisasi client Resend
func getResendClient() *resend.Client {
	return resend.NewClient(os.Getenv("RE_API_KEY"))
}

func SendVerificationEmail(email string, token string) error {
	client := getResendClient()

	link := fmt.Sprintf("https://tujuan-baru-kamu.com/verify-email?token=%s", token)
	body := fmt.Sprintf(`<h2>Verify Email</h2><p>Klik tombol berikut:</p><a href="%s">Verify Email</a>`, link)

	params := &resend.SendEmailRequest{
		From:    "onboarding@resend.dev", // Atau domain terverifikasi kamu
		To:      []string{email},
		Subject: "Verify Your Email",
		Html:    body,
	}

	_, err := client.Emails.Send(params)
	return err
}

func SendResetPasswordEmail(email string, token string) error {
	client := getResendClient()

	link := fmt.Sprintf("https://tujuan-baru-kamu.com/reset-password?token=%s", token)
	body := fmt.Sprintf(`<h2>Reset Password</h2><p>Klik tautan berikut:</p><a href="%s">Reset Password</a>`, link)

	params := &resend.SendEmailRequest{
		From:    "onboarding@resend.dev",
		To:      []string{email},
		Subject: "Permintaan Reset Password",
		Html:    body,
	}

	_, err := client.Emails.Send(params)
	return err
}