package tool

import (
	"github.com/wneessen/go-mail"
)

func SendEmail(toEmail string) error {
	// create the core of the message
	message := mail.NewMsg()
	if err := message.From("atlas@example.com"); err != nil {
		return err
	}
	if err := message.To(toEmail); err != nil {
		return err
	}
	message.Subject("Atlas: verification code")
	message.SetBodyString(mail.TypeTextPlain, "This is your verification code :\n0 0 0 0 0 0\ngo to this url : ....\n")

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