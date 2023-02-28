package mail

import (
	"fmt"
	"kedai/backend/be-kedai/config"

	"gopkg.in/gomail.v2"
)

type MailUtils interface {
	SendUpdatePasswordEmail(receiverEmail string, verificationCode string) error
}

type mailUtilsImpl struct {
	dialer *gomail.Dialer
}

type MailUtilsConfig struct {
	Dialer *gomail.Dialer
}

func NewMailUtils(cfg *MailUtilsConfig) MailUtils {
	return &mailUtilsImpl{
		dialer: cfg.Dialer,
	}
}

func (u *mailUtilsImpl) SendUpdatePasswordEmail(receiverEmail string, verificationCode string) error {
	msg := gomail.NewMessage()
	msg.SetHeader("From", config.GetEnv("MAILER_USERNAME", ""))
	msg.SetHeader("To", receiverEmail)
	msg.SetHeader("Subject", "Update Password")
	msg.SetBody("text/html", fmt.Sprintf("Your verification code is: %s", verificationCode))

	err := u.dialer.DialAndSend(msg)
	if err != nil {
		return err
	}

	return nil
}
