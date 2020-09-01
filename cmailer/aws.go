package cmailer

import (
	"context"

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

func (m *AWSMailer) Send(ctx context.Context, p SendParams) error {
	input := &ses.SendEmailInput{
		Source: aws.String(p.From),
		Destination: &ses.Destination{
			ToAddresses: []*string{&p.To},
		},
		Message: &ses.Message{
			Subject: &ses.Content{
				Charset: aws.String(charset),
				Data:    aws.String(p.Subject),
			},
			Body: &ses.Body{},
		},
	}

	if p.HTMLBody != nil {
		input.Message.Body.Html = &ses.Content{
			Charset: aws.String(charset),
			Data:    p.HTMLBody,
		}
	}

	if p.PlainBody != nil {
		input.Message.Body.Text = &ses.Content{
			Charset: aws.String(charset),
			Data:    p.PlainBody,
		}
	}

	_, err := m.sess.SendEmailWithContext(ctx, input)
	if err != nil {
		return cerror.New(err, "failed to send email", map[string]interface{}{
			"from":      p.From,
			"to":        p.To,
			"subject":   p.Subject,
			"plainBody": p.PlainBody,
			"htmlBody":  p.HTMLBody,
		})
	}

	return nil
}
