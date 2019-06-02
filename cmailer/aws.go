package cmailer

import (
	"github.com/tusharsoni/copper/cerror"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ses"
)

const charset = "UTF-8"

type AWSMailer struct {
	sess *ses.SES
}

func NewAWSMailer(config AWSConfig) (Mailer, error) {
	sess, err := session.NewSession(&aws.Config{
		Region: aws.String(config.Region),
		Credentials: credentials.NewStaticCredentials(
			config.AccessKeyId,
			config.SecretAccessKey,
			"",
		),
	})
	if err != nil {
		return nil, cerror.New(err, "failed to create new aws session", nil)
	}

	return &AWSMailer{
		sess: ses.New(sess),
	}, nil
}

func (m *AWSMailer) SendPlain(from, to, subject, body string) (confirmation string, err error) {
	return m.send(from, to, subject, body, false)
}

func (m *AWSMailer) SendHTML(from, to, subject, body string) (confirmation string, err error) {
	return m.send(from, to, subject, body, true)
}

func (m *AWSMailer) send(from, to, subject, body string, html bool) (confirmation string, err error) {
	input := &ses.SendEmailInput{
		Source: aws.String(from),
		Destination: &ses.Destination{
			ToAddresses: []*string{&to},
		},
		Message: &ses.Message{
			Subject: &ses.Content{
				Charset: aws.String(charset),
				Data:    aws.String(subject),
			},
			Body: &ses.Body{},
		},
	}

	content := &ses.Content{
		Charset: aws.String(charset),
		Data:    aws.String(body),
	}

	if html {
		input.Message.Body.Html = content
	} else {
		input.Message.Body.Text = content
	}

	result, err := m.sess.SendEmail(input)
	if err != nil {
		return "", cerror.New(err, "failed to send email", map[string]interface{}{
			"from":    from,
			"to":      to,
			"subject": subject,
			"body":    body,
		})
	}

	return *result.MessageId, nil
}
