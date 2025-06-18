package models

type Service interface {
	SendWelcomeEmail(to string) error
}
