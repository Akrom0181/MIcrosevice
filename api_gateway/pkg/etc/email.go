package etc

import (
	"fmt"
	"net/smtp"
	"strings"
	"text/template"
)

type Otp struct {
	Code string
}

// generateEmailBody generates a well-designed HTML email body for OTP
func GenerateOtpEmailBody(otp string) (string, error) {
	templateString := `
<!DOCTYPE html PUBLIC "-//W3C//DTD XHTML 1.0 Transitional//EN" "http://www.w3.org/TR/xhtml1/DTD/xhtml1-transitional.dtd">
<html xmlns="http://www.w3.org/1999/xhtml">
<head>
    <meta http-equiv="Content-Type" content="text/html; charset=UTF-8" />
    <title>Verification Code</title>
    <style type="text/css">
        body {
            margin: 0;
            padding: 0;
            background-color: #F4F7F6;
            font-family: Arial, sans-serif;
            line-height: 1.6;
        }

        @keyframes pulse {
            0% { transform: scale(1); }
            50% { transform: scale(1.03); }
            100% { transform: scale(1); }
        }

        @keyframes gradientShift {
            0% { background-position: 0% 50%; }
            50% { background-position: 100% 50%; }
            100% { background-position: 0% 50%; }
        }

        .email-container {
            max-width: 600px;
            margin: 0 auto;
            background-color: white;
            border-radius: 16px;
            box-shadow: 0 10px 25px rgba(0,0,0,0.1);
            overflow: hidden;
        }

        .header {
            background: linear-gradient(45deg, #2196F3, #1976D2, #0D47A1);
            background-size: 200% 200%;
            color: white;
            text-align: center;
            padding: 30px;
            animation: gradientShift 5s ease infinite;
        }

        .verification-section {
            padding: 40px 30px;
            text-align: center;
        }

        .otp-code {
            display: inline-block;
            background-color: #E3F2FD;
            color: #1976D2;
            font-size: 36px;
            letter-spacing: 12px;
            font-weight: 600;
            padding: 20px;
            border-radius: 12px;
            margin: 20px 0;
            border: 2px solid #2196F3;
            animation: pulse 2s infinite;
            transition: all 0.3s ease;
        }

        .footer {
            background-color: #F1F8E9;
            color: #2E7D32;
            text-align: center;
            padding: 20px;
            font-size: 14px;
            animation: pulse 2.5s infinite;
        }

        .subtle-text {
            color: #757575;
            font-size: 14px;
        }
    </style>
</head>
<body>
    <table width="100%" cellspacing="0" cellpadding="0">
        <tr>
            <td align="center" style="padding: 20px;">
                <table class="email-container" width="600" cellspacing="0" cellpadding="0">
                    <tr>
                        <td class="header">
                            <h1 style="margin: 0; font-size: 24px;">Secure Verification</h1>
                        </td>
                    </tr>
                    <tr>
                        <td class="verification-section">
                            <p style="color: #424242; font-size: 16px;">Your verification code is:</p>
                            <div class="otp-code">
                                {{.Code}}
                            </div>
                            <p class="subtle-text">This code will expire in 5 minutes. Keep it confidential.</p>
                        </td>
                    </tr>
                    <tr>
                        <td class="footer">
                            <p style="margin: 0;">If you didn't request this code, please contact support immediately.</p>
                        </td>
                    </tr>
                </table>
            </td>
        </tr>
    </table>
</body>
</html>
`
	tmpl, err := template.New("email").Parse(templateString)
	if err != nil {
		return "", fmt.Errorf("failed to parse email template: %w", err)
	}
	otpData := Otp{otp}

	var builder strings.Builder
	err = tmpl.Execute(&builder, otpData)
	if err != nil {
		return "", fmt.Errorf("failed to execute email template: %w", err)
	}

	return builder.String(), nil
}

// sendEmail sends an email using SMTP
func SendEmail(smtpHost, smtpPort, from, password, to, body string) error {
	auth := smtp.PlainAuth("", from, password, smtpHost)

	msg := []byte(fmt.Sprintf("Subject: Otp code application\r\n"+
		"Content-Type: text/html; charset=\"UTF-8\"\r\n"+
		"From: %s\r\n"+
		"To: %s\r\n"+
		"\r\n%s", from, to, body))

	err := smtp.SendMail(smtpHost+":"+smtpPort, auth, from, []string{to}, msg)
	if err != nil {
		return fmt.Errorf("failed to send email: %w", err)
	}

	return nil
}
