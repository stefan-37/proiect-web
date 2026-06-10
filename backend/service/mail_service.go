package service

import(
	"net/smtp"
	"os"
	"log"
)

func sendEmail(to string)error{
	from := os.Getenv("EMAIL")
	password := os.Getenv("EMAIL_PASSWORD")

	if from == "" || password == "" {
		log.Printf("Email credentials not set in environment variables")
		return nil
	}

	smtpHost := "smtp.gmail.com"
	smtpPort := "587"

	auth := smtp.PlainAuth("", from, password, smtpHost)

	message := []byte("To: " + to + "\r\n" +
      "Subject: Hello from Backend\r\n" +
      "\r\n" +
      "This is a verification email.\r\n")

	return smtp.SendMail(smtpHost+":"+smtpPort, auth, from, []string{to}, message)
}

func sendPasswordEmail(to string, newPassword string) error {
	from := os.Getenv("EMAIL")
	password := os.Getenv("EMAIL_PASSWORD")

	if from == "" || password == "" {
		log.Printf("Email credentials not set in environment variables")
		return nil
	}

	smtpHost := "smtp.gmail.com"
	smtpPort := "587"

	auth := smtp.PlainAuth("", from, password, smtpHost)

	message := []byte("To: " + to + "\r\n" +
      "Subject: Password Reset\r\n" +
      "\r\n" +
      "New password: " + newPassword + "\r\n")

	return smtp.SendMail(smtpHost+":"+smtpPort, auth, from, []string{to}, message)
}