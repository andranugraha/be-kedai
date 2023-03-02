package connection

import (
	"kedai/backend/be-kedai/config"

	"github.com/mailgun/mailgun-go/v4"
)

var (
	mailgunConfig = config.Mailgun
	mailer        *mailgun.MailgunImpl
)

func ConnectMailer() {
	mailer = mailgun.NewMailgun(mailgunConfig.DOMAIN, mailgunConfig.PRIVATE_API_KEY)
}

func GetMailer() *mailgun.MailgunImpl {
	return mailer
}
