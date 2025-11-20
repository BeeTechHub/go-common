package awsSesV2

import (
	"errors"
	"sync"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials/stscreds"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ses"
	"github.com/aws/aws-sdk-go/service/sts"
)

const CharSet = "UTF-8"

type SesV2Wrapper struct {
	SesV2       *ses.SES
	EmailSender string
}

type SesAccountConfig struct {
	RoleARN     string
	Region      string // Default region is ap-southeast-1
	EmailSender string
}

// SesRouter quản lý routing email dựa trên system identifier (GS, TR, ...)
type SesRouter struct {
	systemMap map[string]SesV2Wrapper // Map system identifier -> SES wrapper
	mu        sync.RWMutex
}

// NewSesRouter tạo một SesRouter mới
func NewSesRouter() *SesRouter {
	return &SesRouter{
		systemMap: make(map[string]SesV2Wrapper),
	}
}

// RegisterSystem đăng ký một hệ thống với AWS SES account tương ứng
// systemID: định danh hệ thống (ví dụ: "GS", "TR")
// account: cấu hình AWS SES account cho hệ thống đó
func (r *SesRouter) RegisterSystem(systemID string, account SesAccountConfig) error {
	if systemID == "" {
		return errors.New("systemID cannot be empty")
	}

	wrapper, err := InitSesWithCredentials(account)
	if err != nil {
		return err
	}

	r.mu.Lock()
	defer r.mu.Unlock()
	r.systemMap[systemID] = wrapper

	return nil
}

// SendEmail gửi email thông qua AWS SES tương ứng với systemID
// systemID: định danh hệ thống (ví dụ: "GS", "TR")
// recipient: địa chỉ email người nhận
// subject: tiêu đề email
// body: nội dung email
func (r *SesRouter) SendEmail(systemID string, recipient string, subject string, body string) (*ses.SendEmailOutput, error) {
	if systemID == "" {
		return nil, errors.New("systemID cannot be empty")
	}

	r.mu.RLock()
	wrapper, exists := r.systemMap[systemID]
	r.mu.RUnlock()

	if !exists {
		return nil, errors.New("system not found: " + systemID)
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
		Source: aws.String(wrapper.EmailSender),
		// Uncomment to use a configuration set
		//ConfigurationSetName: aws.String(ConfigurationSet),
	}

	// return wrapper.SesV2.SendEmail(recipient, subject, body)
	return wrapper.SesV2.SendEmail(input)
}

// NewSesRouterWithSystems tạo một SesRouter mới và đăng ký nhiều systems cùng lúc
// systems: map systemID -> SesAccountConfig
func NewSesRouterWithSystems(systems map[string]SesAccountConfig) (*SesRouter, error) {
	router := NewSesRouter()

	for systemID, account := range systems {
		if err := router.RegisterSystem(systemID, account); err != nil {
			return nil, err
		}
	}

	return router, nil
}

// InitSesWithCredentials khởi tạo SES với IAM role ARN
func InitSesWithCredentials(config SesAccountConfig) (SesV2Wrapper, error) {
	if config.RoleARN == "" {
		return SesV2Wrapper{}, errors.New("roleARN is required")
	}

	if config.Region == "" {
		return SesV2Wrapper{}, errors.New("region is required")
	}

	if config.EmailSender == "" {
		return SesV2Wrapper{}, errors.New("emailSender is required")
	}

	// Tạo base session sử dụng default credentials chain (ECS task role, instance profile, etc.)
	baseSess, err := session.NewSession(&aws.Config{
		Region: aws.String(config.Region),
	})
	if err != nil {
		return SesV2Wrapper{}, err
	}

	// Assume role để lấy credentials từ role ARN
	stsClient := sts.New(baseSess)
	creds := stscreds.NewCredentialsWithClient(stsClient, config.RoleARN)

	// Tạo session với assumed role credentials
	sess, err := session.NewSession(&aws.Config{
		Region:      aws.String(config.Region),
		Credentials: creds,
	})
	if err != nil {
		return SesV2Wrapper{}, err
	}

	sesClient := ses.New(sess)
	return SesV2Wrapper{sesClient, config.EmailSender}, nil
}
