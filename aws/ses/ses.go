package awsSes

import (
	"bytes"
	"encoding/base64"
	"errors"
	"fmt"
	"html"
	"mime/multipart"
	"net/textproto"
	"strings"

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

func (sesWrapper SesWrapper) SendRichEmail(recipients []string, subject string, textBody string, htmlBody string, attachments []EmailAttachment,
) (*ses.SendRawEmailOutput, error) {
	if sesWrapper.Ses == nil {
		return nil, nilSesError
	}

	if textBody == "" && htmlBody == "" {
		return nil, errors.New("email body must not be empty")
	}

	/*if textBody == "" {
		if _textBody, err := html2text.FromString(htmlBody); err == nil {
			textBody = _textBody
		}
	}*/

	if htmlBody == "" {
		htmlBody = "<pre>" + html.EscapeString(textBody) + "</pre>"
	}

	var raw bytes.Buffer
	mixedWriter := multipart.NewWriter(&raw)

	// ===== Headers =====
	headers := map[string]string{
		"From":         sesWrapper.EmailSender,
		"To":           strings.Join(recipients, ","),
		"Subject":      subject,
		"MIME-Version": "1.0",
		"Content-Type": fmt.Sprintf("multipart/mixed; boundary=%s", mixedWriter.Boundary()),
	}

	for k, v := range headers {
		raw.WriteString(fmt.Sprintf("%s: %s\r\n", k, v))
	}
	raw.WriteString("\r\n")

	// ===== multipart/alternative =====
	altBuffer := bytes.Buffer{}
	altWriter := multipart.NewWriter(&altBuffer)

	// Text part
	textPart, err := altWriter.CreatePart(
		textproto.MIMEHeader{
			"Content-Type": {"text/plain; charset=UTF-8"},
		},
	)
	if err != nil {
		return nil, err
	}
	_, _ = textPart.Write([]byte(textBody))

	// HTML part
	htmlPart, err := altWriter.CreatePart(
		textproto.MIMEHeader{
			"Content-Type": {"text/html; charset=UTF-8"},
		},
	)
	if err != nil {
		return nil, err
	}
	_, _ = htmlPart.Write([]byte(htmlBody))

	_ = altWriter.Close()

	// Gắn alternative vào mixed
	altPart, err := mixedWriter.CreatePart(
		textproto.MIMEHeader{
			"Content-Type": {fmt.Sprintf("multipart/alternative; boundary=%s", altWriter.Boundary())},
		},
	)
	if err != nil {
		return nil, err
	}
	_, _ = altPart.Write(altBuffer.Bytes())

	// ===== Attachments =====
	for _, att := range attachments {
		part, err := mixedWriter.CreatePart(
			textproto.MIMEHeader{
				"Content-Type":              {string(att.ContentType)},
				"Content-Disposition":       {fmt.Sprintf(`attachment; filename="%s"`, att.Filename)},
				"Content-Transfer-Encoding": {"base64"},
			},
		)
		if err != nil {
			return nil, err
		}

		encoder := base64.NewEncoder(base64.StdEncoding, part)
		_, _ = encoder.Write(att.Data)
		_ = encoder.Close()
	}

	_ = mixedWriter.Close()

	// ===== Send via SES =====
	input := &ses.SendRawEmailInput{
		RawMessage: &ses.RawMessage{
			Data: raw.Bytes(),
		},
	}

	return sesWrapper.Ses.SendRawEmail(input)
}
