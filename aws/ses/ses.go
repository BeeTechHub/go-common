package awsSes

import (
	"errors"

	config "github.com/BeeTechHub/go-common/aws/config"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ses"
)

var nilSesError = errors.New("Access ses failed because ses nil")

const CharSet = "UTF-8"

var svc *ses.SES

type SesWrapper struct {
	Ses         *ses.SES
	EmailSender string
}

func InitSes(emailSender string) SesWrapper {
	if svc == nil {
		sess := config.GetAWSSession()
		svc = ses.New(sess)
	}

	return SesWrapper{svc, emailSender}
}

func (sesWrapper SesWrapper) SendEmail(recipient string, subject string, body string) (*ses.SendEmailOutput, error) {
	if sesWrapper.Ses == nil {
		return nil, nilSesError
	}

	// Assemble the email.
	input := &ses.SendEmailInput{
		Destination: &ses.Destination{
			CcAddresses: []*string{},
			ToAddresses: []*string{
				aws.String(recipient),
			},
		},
		Message: &ses.Message{
			Body: &ses.Body{
				Html: &ses.Content{
					Charset: aws.String(CharSet),
					Data:    aws.String(body),
				},
				Text: &ses.Content{
					Charset: aws.String(CharSet),
					Data:    aws.String(body),
				},
			},
			Subject: &ses.Content{
				Charset: aws.String(CharSet),
				Data:    aws.String(subject),
			},
		},
		Source: aws.String(sesWrapper.EmailSender),
		// Uncomment to use a configuration set
		//ConfigurationSetName: aws.String(ConfigurationSet),
	}

	// Attempt to send the email.
	result, err := sesWrapper.Ses.SendEmail(input)

	// Display error messages if they occur.
	if err != nil {
		return nil, err
	}

	return result, nil
}
