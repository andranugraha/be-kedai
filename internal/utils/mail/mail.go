package mail

import (
	"context"
	"kedai/backend/be-kedai/config"
	"time"

	"github.com/mailgun/mailgun-go/v4"
)

type MailUtils interface {
	SendUpdatePasswordEmail(receiverEmail string, verificationCode string) error
}

type mailUtilsImpl struct {
	mailer *mailgun.MailgunImpl
}

type MailUtilsConfig struct {
	Mailer *mailgun.MailgunImpl
}

func NewMailUtils(cfg *MailUtilsConfig) MailUtils {
	return &mailUtilsImpl{
		mailer: cfg.Mailer,
	}
}

func (u *mailUtilsImpl) SendUpdatePasswordEmail(receiverEmail string, verificationCode string) error {
	sender := config.GetEnv("MAILGUN_SENDER", "Support@kedai.com")
	subject := "Update Password Verification Code"
	body := "Your verification code is: " + verificationCode

	msg := u.mailer.NewMessage(sender, subject, body, receiverEmail)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	_, _, err := u.mailer.Send(ctx, msg)

	if err != nil {
		return err
	}

	return nil
}
