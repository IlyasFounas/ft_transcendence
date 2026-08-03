package tool

import (
	"os"
	"fmt"
	"github.com/wneessen/go-mail"
	"github.com/resend/resend-go/v3"
)

func SendEmailProd(toEmail string, subject string, html string) error {
    client := resend.NewClient(os.Getenv("RESEND_API_KEY"))

    params := &resend.SendEmailRequest{
        From:    "Acme <atlas-no-reply@ocefunding.com>",
        To:      []string{toEmail},
        Html:    html,
        Subject: subject,
    }

    _, err := client.Emails.Send(params)
    if err != nil {
        return err
    }
	return nil
}

func SendEmail(toEmail string, message *mail.Msg) error {
	// send the email with the new client
	client, err := mail.NewClient("mailpit",
		mail.WithPort(1025),
		mail.WithTLSPolicy(mail.NoTLS),
		mail.WithSMTPAuth(mail.SMTPAuthNoAuth),
	)
	if err != nil {
		return err
	}
	if err := client.DialAndSend(message); err != nil {
		return err
	}
	return nil
}

func returnActivationLink(activationJWT string) string {
    return fmt.Sprintf(
        "This is your activation link: http://localhost:5130/api/user/activation?token=%s\n%s",
        activationJWT,
		"you have 24 hours to activate your account before it gets deleted",
	)
}

func SendActivationEmail(toEmail string, activationJWT string) error {
	// create the message

	env := os.Getenv("ENV")
	subject := "Atlas: verification code"
	html := returnActivationLink(activationJWT)
	if len(env) != 0 && env == "prod" {
		if err := SendEmailProd(toEmail, subject, html); err != nil {
			return err
		}
	} else {
		message := mail.NewMsg()
		if err := message.From("atlas@example.com"); err != nil {
			return err
		}
		if err := message.To(toEmail); err != nil {
			return err
		}
		message.Subject(subject)
		message.SetBodyString(mail.TypeTextPlain, html)
		if err := SendEmail(toEmail, message); err != nil {
			return err
		}
	}
	return nil
}