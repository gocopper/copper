package ctexter

type Svc interface {
	SendSMS(phoneNumber, message string) (confirmationCode string, err error)
}
