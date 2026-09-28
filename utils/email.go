// File: utils/email.go
package utils

import (
	"crypto/rand"
	"math/big"
	"os"

	"gopkg.in/gomail.v2"
)

// GenerateOTP membuat 6 digit angka acak
func GenerateOTP() string {
	const letters = "0123456789"
	b := make([]byte, 6)
	for i := range b {
		n, _ := rand.Int(rand.Reader, big.NewInt(int64(len(letters))))
		b[i] = letters[n.Int64()]
	}
	return string(b)
}

// SendOTPEmail mengirim email via SMTP Gmail
func SendOTPEmail(toEmail, otpCode string) error {
	// Pastikan Anda mengatur env ini di Vercel atau .env lokal (Pakai App Password Gmail)
	smtpHost := os.Getenv("SMTP_HOST") // smtp.gmail.com
	smtpPort := 587                    // 587
	smtpUser := os.Getenv("SMTP_USER") // email.resto@gmail.com
	smtpPass := os.Getenv("SMTP_PASS") // App Password 16 digit

	m := gomail.NewMessage()
	m.SetHeader("From", smtpUser)
	m.SetHeader("To", toEmail)
	m.SetHeader("Subject", "Kode OTP Reset Password - Kecilung Resto")
	m.SetBody("text/html", "<h3>Kode OTP Anda: <b>"+otpCode+"</b></h3><p>Kode ini hanya berlaku selama 10 menit. Jangan berikan kode ini kepada siapapun.</p>")

	d := gomail.NewDialer(smtpHost, smtpPort, smtpUser, smtpPass)
	return d.DialAndSend(m)
}