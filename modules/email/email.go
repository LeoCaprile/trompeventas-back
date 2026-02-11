package email

import (
	"fmt"
	"os"
	"strings"

	"github.com/charmbracelet/log"
	"github.com/resend/resend-go/v3"
)

var EmailClient *resend.Client

func InitResendClient() {
	apiKey := os.Getenv("RESEND_API_KEY")
	if apiKey == "" {
		log.Warn("RESEND_API_KEY environment variable is not set - email functionality will not work")
		return
	}

	EmailClient = resend.NewClient(apiKey)
	log.Infof("Resend client initialized with API key: %s...%s", apiKey[:7], apiKey[len(apiKey)-4:])
}

func SendVerificationEmail(toEmail, userName, verificationLink string) error {
	// Check if EmailClient is initialized
	if EmailClient == nil {
		return fmt.Errorf("email client not initialized - check RESEND_API_KEY env variable")
	}

	// Read template
	templateBytes, err := os.ReadFile("modules/email/templates/verification_email.html")
	if err != nil {
		return fmt.Errorf("failed to read email template: %w", err)
	}

	htmlContent := string(templateBytes)
	htmlContent = strings.ReplaceAll(htmlContent, "{{USER_NAME}}", userName)
	htmlContent = strings.ReplaceAll(htmlContent, "{{VERIFICATION_LINK}}", verificationLink)

	fromEmail := os.Getenv("EMAIL_FROM")
	if fromEmail == "" {
		fromEmail = "Trompeventas <contacto@trompeventas.cl>"
	}

	params := &resend.SendEmailRequest{
		From:    fromEmail,
		To:      []string{toEmail},
		Subject: "Verifica tu Email - Trompeventas",
		Html:    htmlContent,
	}

	log.Infof("Sending verification email to: %s", toEmail)
	log.Debugf("Email from: %s | Subject: %s", params.From, params.Subject)

	sent, err := EmailClient.Emails.Send(params)
	if err != nil {
		log.Error("Failed to send verification email", "error", err, "to", toEmail)
		return fmt.Errorf("failed to send verification email: %w", err)
	}

	log.Infof("Verification email sent successfully! ID: %s", sent.Id)
	return nil
}
