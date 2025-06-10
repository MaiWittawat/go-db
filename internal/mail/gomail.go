package mail

type Mail interface {
	SendEmail(msg string, subject string, to []string) error
	SendWelcomeEmail(to []string) error 
}

