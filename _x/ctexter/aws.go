package ctexter

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sns"
	"github.com/tusharsoni/copper/cerror"
)

type awsSvc struct {
	snsClient *sns.SNS
}

func newAWSSvc(config AWSConfig) (Svc, error) {
	sess, err := session.NewSession(&aws.Config{
		Region: aws.String(config.Region),
		Credentials: credentials.NewStaticCredentials(
			config.AccessKeyId,
			config.SecretAccessKey,
			"",
		),
	})
	if err != nil {
		return nil, cerror.New(err, "failed to create new AWS session", nil)
	}

	return &awsSvc{
		snsClient: sns.New(sess),
	}, nil
}

func (s *awsSvc) SendSMS(phoneNumber, message string) (confirmationCode string, err error) {
	publishInput := &sns.PublishInput{
		PhoneNumber: aws.String(phoneNumber),
		Message:     aws.String(message),
	}

	publishOutput, err := s.snsClient.Publish(publishInput)
	if err != nil {
		return "", err
	}

	return *publishOutput.MessageId, nil
}
