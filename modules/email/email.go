package email

import (
	"fmt"
	"os"
	"strings"

	"github.com/resend/resend-go/v3"
)

var EmailClient *resend.Client

func InitResendClient() {
	EmailClient = resend.NewClient(os.Getenv("RESEND_API_KEY"))
}

func SendVerificationEmail(toEmail, userName, verificationLink string) error {
	// Read template
	templateBytes, err := os.ReadFile("modules/email/templates/verification_email.html")
	if err != nil {
		return fmt.Errorf("failed to read email template: %w", err)
	}

	htmlContent := string(templateBytes)
	htmlContent = strings.ReplaceAll(htmlContent, "{{USER_NAME}}", userName)
	htmlContent = strings.ReplaceAll(htmlContent, "{{VERIFICATION_LINK}}", verificationLink)

	params := &resend.SendEmailRequest{
		From:    "Trompeventas <noreply@trompeventas.cl>",
		To:      []string{toEmail},
		Subject: "Verifica tu Email - Trompeventas",
		Html:    htmlContent,
	}

	_, err = EmailClient.Emails.Send(params)
	if err != nil {
		return fmt.Errorf("failed to send verification email: %w", err)
	}

	return nil
}
