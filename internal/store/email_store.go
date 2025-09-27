package store

import (
	"crypto/rand"
	"math/big"
	"os"

	"github.com/resend/resend-go/v2"
)

type ResendEmailBody struct {
	From    string   `json:"from"`
	To      []string `json:"to"`
	Subject string   `json:"subject"`
	HTML    string   `json:"html"`
}

type EmailService interface {
	SendResetCode(body ResendEmailBody) error
}

type ResendEmailService struct {
	client *resend.Client
}

func NewResendEmailService() *ResendEmailService {
	apiKey := os.Getenv("RESEND_API_KEY")
	client := resend.NewClient(apiKey)
	return &ResendEmailService{

		client: client,
	}
}

func (s *ResendEmailService) SendResetCode(body ResendEmailBody) error {
	params := &resend.SendEmailRequest{
		From:    body.From,
		To:      body.To,
		Subject: body.Subject,
		Html:    body.HTML,
	}
	_, err := s.client.Emails.Send(params)
	if err != nil {
		return err
	}
	return nil
}

func (s *ResendEmailService) GenerateResetEmailHTML(code string) string {
	return `
    <div style="font-family: Arial, sans-serif; max-width: 600px; margin: 0 auto;">
        <h2>Password Reset Request</h2>
        <p>You requested a password reset. Use the code below to reset your password:</p>
        <div style="background-color: #f0f0f0; padding: 20px; text-align: center; font-size: 24px; font-weight: bold; margin: 20px 0;">
			` + code + `
        </div>
        <p>This code will expire in 15 minutes.</p>
        <p>If you didn't request this, please ignore this email.</p>
    </div>
    `
}

func GenerateResetCode() (string, error) {
	const digits = "0123456789"
	code := make([]byte, 6)

	for i := range code {
		num, err := rand.Int(rand.Reader, big.NewInt(int64(len(digits))))
		if err != nil {
			return "", err
		}
		code[i] = digits[num.Int64()]
	}
	return string(code), nil
}
