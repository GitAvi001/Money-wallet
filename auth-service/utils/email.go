package utils

import (
	"auth-service/config"
	"fmt"
	"log"
	"net/smtp" //SMTP server package for Go
)

func SendVerificationEmail(toEmail, token string) {
	//using .env configurations for SMTP server
	smtpHost := config.GetEnv("SMTP_HOST", "smtp.gmail.com")
	smtpPort := config.GetEnv("SMTP_PORT", "587")
	smtpUsername := config.GetEnv("SMTP_USERNAME", "")
	smtpPassword := config.GetEnv("SMTP_PASSWORD", "")
	fromEmail := config.GetEnv("FROM_EMAIL", "")
	frontendURL := config.GetEnv("FRONTEND_URL", "http://localhost:5173")

	if smtpUsername == "" || smtpPassword == "" {
		log.Println("SMTP credentials not configured. Verification email not sent.")
		log.Printf("Verification link: %s/verify-email?token=%s", frontendURL, token)
		return
	}

	verificationURL := fmt.Sprintf("%s/verify-email?token=%s", frontendURL, token)

	//HTML body showing in the email content
	subject := "Email Verification - Money Transfer App"
	body := fmt.Sprintf(`
<!DOCTYPE html>
<html>
<head>
    <style>
        body { font-family: Arial, sans-serif; line-height: 1.6; color: #333; }
        .container { max-width: 600px; margin: 0 auto; padding: 20px; }
        .header { background-color: #4F46E5; color: white; padding: 20px; text-align: center; }
        .content { background-color: #f9f9f9; padding: 30px; text-align: center; }
        .button { 
            background-color: #4F46E5 !important; 
            color: #ffffff !important; 
            padding: 15px 40px; 
            text-decoration: none; 
            border-radius: 5px; 
            display: inline-block; 
            margin: 20px 0; 
            font-weight: bold;
            font-size: 16px;
        }
        .footer { text-align: center; margin-top: 20px; color: #666; font-size: 12px; }
    </style>
</head>
<body>
    <div class="container">
        <div class="header">
            <h1>Verify Your Email</h1>
        </div>
        <div class="content">
            <p>Thank you for registering with our Money Transfer App!</p>
            <p>Please click the button below to verify your email address:</p>
            <a href="%s" class="button" style="background-color: #4F46E5; color: #ffffff; padding: 15px 40px; text-decoration: none; border-radius: 5px; display: inline-block; font-weight: bold; font-size: 16px;">Verify Email</a>
            <p>Or copy and paste this link into your browser:</p>
            <p style="word-break: break-all; color: #4F46E5;">%s</p>
            <p>This link will expire in 24 hours.</p>
            <p>If you didn't create an account, please ignore this email.</p>
        </div>
        <div class="footer">
            <p>© 2025 Money Transfer App. All rights reserved.</p>
        </div>
    </div>
</body>
</html>
	`, verificationURL, verificationURL)

	//message for email content
	message := []byte(fmt.Sprintf("To: %s\r\n"+
		"Subject: %s\r\n"+
		"MIME-Version: 1.0\r\n"+
		"Content-Type: text/html; charset=UTF-8\r\n"+
		"\r\n"+
		"%s\r\n", toEmail, subject, body))

	auth := smtp.PlainAuth("", smtpUsername, smtpPassword, smtpHost)
	err := smtp.SendMail(smtpHost+":"+smtpPort, auth, fromEmail, []string{toEmail}, message)

	if err != nil {
		log.Printf("Failed to send email to %s: %v", toEmail, err)
		log.Printf("Verification link: %s", verificationURL)
	} else {
		log.Printf("Verification email sent to %s", toEmail)
	}
}
