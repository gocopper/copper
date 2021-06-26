package cmailer

import (
	"context"

	"github.com/aws/aws-sdk-go/aws"
	awscredentials "github.com/aws/aws-sdk-go/aws/credentials"
	awssession "github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ses"
	"github.com/gocopper/copper/cconfig"
	"github.com/gocopper/copper/cerrors"
)

const charsetUTF8 = "UTF-8"

// AWSConfig is used to configure the AWS mailer
type AWSConfig struct {
	Region          string `toml:"region"`
	AccessKeyID     string `toml:"access_key_id"`
	SecretAccessKey string `toml:"secret_access_key"`
}

// NewAWSMailer creates an implementation of Mailer that uses AWS
func NewAWSMailer(appConfig cconfig.Config) (Mailer, error) {
	var config AWSConfig

	err := appConfig.Load("aws", &config)
	if err != nil {
		return nil, cerrors.New(err, "failed to load AWS config", nil)
	}

	sess, err := awssession.NewSession(&aws.Config{
		Region: aws.String(config.Region),
		Credentials: awscredentials.NewStaticCredentials(
			config.AccessKeyID,
			config.SecretAccessKey,
			"",
		),
	})
	if err != nil {
		return nil, cerrors.New(err, "failed to create new AWS session", nil)
	}

	return &awsMailer{
		sess: ses.New(sess),
	}, nil
}

type awsMailer struct {
	sess *ses.SES
}

func (m *awsMailer) Send(ctx context.Context, p SendParams) error {
	input := &ses.SendEmailInput{
		Source: aws.String(p.From),
		Destination: &ses.Destination{
			ToAddresses: make([]*string, len(p.To)),
		},
		Message: &ses.Message{
			Subject: &ses.Content{
				Charset: aws.String(charsetUTF8),
				Data:    aws.String(p.Subject),
			},
			Body: &ses.Body{},
		},
	}

	for i := range p.To {
		input.Destination.ToAddresses[i] = &p.To[i]
	}

	if p.HTMLBody != nil {
		input.Message.Body.Html = &ses.Content{
			Charset: aws.String(charsetUTF8),
			Data:    p.HTMLBody,
		}
	}

	if p.PlainBody != nil {
		input.Message.Body.Text = &ses.Content{
			Charset: aws.String(charsetUTF8),
			Data:    p.PlainBody,
		}
	}

	_, err := m.sess.SendEmailWithContext(ctx, input)
	if err != nil {
		return cerrors.New(err, "failed to send email", map[string]interface{}{
			"from":      p.From,
			"to":        p.To,
			"subject":   p.Subject,
			"plainBody": p.PlainBody,
			"htmlBody":  p.HTMLBody,
		})
	}

	return nil
}
