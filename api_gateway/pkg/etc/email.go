package etc

import (
	"fmt"
	"log"
	"net/smtp"
	"strings"
	"text/template"
	"user_api_gateway/config"
)

type Otp struct {
	Code string
}

func GenerateOtpEmailBody(otp string) (string, error) {
	templateString := `
<!DOCTYPE html>
<html>
<body>
    <p>Your OTP to verify your App account: <strong>{{.Code}}</strong></p>
</body>
</html>
`
	tmpl, err := template.New("email").Parse(templateString)
	if err != nil {
		return "", fmt.Errorf("failed to parse email template: %w", err)
	}

	otpData := Otp{Code: otp}

	var builder strings.Builder
	err = tmpl.Execute(&builder, otpData)
	if err != nil {
		return "", fmt.Errorf("failed to execute email template: %w", err)
	}

	return builder.String(), nil
}

// SendMail sends an email with the generated OTP message
func SendMail(toEmail string, msg string) error {
	from := config.SmtpUsername
	to := []string{toEmail}
	subject := "Register your app"

	// Correctly format the email as an HTML message
	body := "MIME-Version: 1.0\r\n" +
		"Content-Type: text/html; charset=\"UTF-8\"\r\n" +
		"From: " + from + "\r\n" +
		"To: " + toEmail + "\r\n" +
		"Subject: " + subject + "\r\n\r\n" +
		msg

	auth := smtp.PlainAuth("", config.SmtpUsername, config.SmtpPassword, config.SmtpServer)

	err := smtp.SendMail(config.SmtpServer+":"+config.SmtpPort, auth, from, to, []byte(body))
	if err != nil {
		log.Println("Error sending mail:", err)
		return err
	}

	log.Println("Email sent successfully to:", toEmail)
	return nil
}
