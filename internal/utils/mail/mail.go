package mail

import (
	"context"
	"kedai/backend/be-kedai/config"
	"time"

	"github.com/mailgun/mailgun-go/v4"
)

type MailUtils interface {
	SendUpdatePasswordEmail(receiverEmail string, verificationCode string) error
	SendUpdatePinEmail(receiverEmail string, verificationCode string) error
	SendResetPasswordEmail(receiverEmail string, token string) error
	SendResetPinEmail(receiverEmail string, token string) error
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

func (u *mailUtilsImpl) SendUpdatePinEmail(receiverEmail string, verificationCode string) error {
	sender := config.GetEnv("MAILGUN_SENDER", "Support@kedai.com")
	subject := "Update Wallet PIN Verification Code"
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

func (u *mailUtilsImpl) SendResetPasswordEmail(receiverEmail string, token string) error {
	sender := config.GetEnv("MAILGUN_SENDER", "Support@kedai.com")
	subject := "Reset Password"
	body := "Please click this link to reset your password: " +
		config.GetEnv("FRONTEND_URL", "http://localhost:3000") +
		config.GetEnv("RESET_PASSWORD_URL", "/reset-password?token=") +
		token

	msg := u.mailer.NewMessage(sender, subject, body, receiverEmail)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	_, _, err := u.mailer.Send(ctx, msg)

	if err != nil {
		return err
	}

	return nil
}

func (u *mailUtilsImpl) SendResetPinEmail(receiverEmail string, token string) error {
	sender := config.GetEnv("MAILGUN_SENDER", "Support@kedai.com")
	subject := "Reset Wallet PIN"
	body := "Please click this link to reset your wallet PIN: " +
		config.GetEnv("FRONTEND_URL", "http://localhost:3000") +
		config.GetEnv("RESET_PASSWORD_URL", "/reset-pin?token=") +
		token

	msg := u.mailer.NewMessage(sender, subject, body, receiverEmail)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	_, _, err := u.mailer.Send(ctx, msg)

	if err != nil {
		return err
	}

	return nil
}
