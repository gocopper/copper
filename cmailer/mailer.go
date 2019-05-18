package cmailer

type Mailer interface {
	SendPlain(from, to, subject, body string) (confirmation string, err error)
	SendHTML(from, to, subject, body string) (confirmation string, err error)
}
