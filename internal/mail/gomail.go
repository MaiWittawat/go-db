package mail

import (
	"crypto/tls"
	"errors"
	appcore_config "go-rebuild/cmd/go-rebuild/config"
	"strconv"

	log "github.com/sirupsen/logrus"
	"gopkg.in/gomail.v2"
)

var (
	ErrSendMessage = errors.New("failed to send message by gomail")
)

type mailService struct {
	mailClient *gomail.Dialer
}

func InitSMTP() *gomail.Dialer {
	port, err := strconv.Atoi(appcore_config.Config.EmailSMTPPort)
	if err != nil {
		log.Panicf("fail to catch type string to int in initSMTP :%v", err)
	}

	mailer := gomail.NewDialer(
		appcore_config.Config.EmailSTMPHost,
		port,
		appcore_config.Config.EmailSMTPUser,
		appcore_config.Config.EmailSMTPPassword,
	)

	mailer.TLSConfig = &tls.Config{InsecureSkipVerify: true}
	return mailer
}

func NewMailService(mailClient *gomail.Dialer) Mail {
	return &mailService{mailClient: mailClient}
}

func (s *mailService) SendEmail(msg string, subject string, to []string) error {
	var baseLogFields = log.Fields{
		"sendTo":   to[0],
		"layer":     "mail_service",
		"operation": "send_welcome_email",
	}

	m := gomail.NewMessage()
	m.SetHeader("From", appcore_config.Config.EmailSMTPFrom)
	m.SetHeader("To", to[0])
	m.SetHeader("Subject", subject)
	m.SetBody("text/html", msg)

	if err := s.mailClient.DialAndSend(m); err != nil {
		log.WithError(err).WithFields(baseLogFields)
		return ErrSendMessage
	}

	log.Info("[gomail]: send email success")
	return nil
}

func (s *mailService) SendWelcomeEmail(to []string) error {
	var baseLogFields = log.Fields{
		"sendTo":   to[0],
		"layer":     "mail_service",
		"operation": "send_welcome_email",
	}

	subject := "Welcome to go-rebuild project"
	msg := "Hello welcome to go-rebuild project"
	m := gomail.NewMessage()
	m.SetHeader("From", appcore_config.Config.EmailSMTPFrom)
	m.SetHeader("To", to[0])
	m.SetHeader("Subject", subject)
	m.SetBody("text/html", msg)

	if err := s.mailClient.DialAndSend(m); err != nil {
		log.WithError(err).WithFields(baseLogFields)
		return ErrSendMessage
	}

	log.Info("[gomail]: send welcome email success")
	return nil
}
