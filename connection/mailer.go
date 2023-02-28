package connection

import (
	"kedai/backend/be-kedai/config"
	"strconv"

	gomail "gopkg.in/gomail.v2"
)

var (
	mailerConfig = config.Mailer
	mailer       *gomail.Dialer
)

func ConnectMailer() {
	portInt, _ := strconv.Atoi(mailerConfig.Port)
	mailer = gomail.NewDialer(mailerConfig.Host, portInt, mailerConfig.Username, mailerConfig.Password)
}

func GetMailer() *gomail.Dialer {
	return mailer
}
