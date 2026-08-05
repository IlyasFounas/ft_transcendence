package mailing

import (
	"fmt"
	"github.com/wneessen/go-mail"
)

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
		"you have 24 hours before your account gets deleted",
	)
}

func SendActivationEmail(toEmail string, activationJWT string) error {
	// create the message
	message := mail.NewMsg()
	if err := message.From("atlas@example.com"); err != nil {
		return err
	}
	if err := message.To(toEmail); err != nil {
		return err
	}
	message.Subject("Atlas: verification code")
	message.SetBodyString(mail.TypeTextPlain, returnActivationLink(activationJWT))
	if err := SendEmail(toEmail, message); err != nil {
		return err
	}
	return nil
}
