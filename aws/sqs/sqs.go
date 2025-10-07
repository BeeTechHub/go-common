package awsSqs

import (
	"errors"
	"fmt"

	config "github.com/BeeTechHub/go-common/aws/config"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/sqs"
)

var nilSqsError = errors.New("Access sqs failed because sqs nil")

var svc *sqs.SQS

type SqsWrapper struct {
	Sqs               *sqs.SQS
	QueueUrl          string
	DelaySeconds      int64
	MaxNumberOfMsg    int64
	WaitTime          int64
	VisibilityTimeout int64
}

// queueName: Tên queue
// delaySeconds: Số giây delay trước khi bản tin có thể pull về
// maxNumberOfMsg: Số lượng bản tin pull về tối đa
// waitTime: Thời gian chờ (giây) tối đa để pull bản tin về
// visibilityTimeout: Thời gian (giây) bản tin ẩn khỏi các subscriber sau khi được pull về (nếu không bị delete)
func InitSqs(queueName string, delaySeconds, maxNumberOfMsg, waitTime, visibilityTimeout int64) (*SqsWrapper, error) {
	if svc == nil {
		sess := config.GetAWSSession()
		svc = sqs.New(sess)
	}

	queueUrlOutput, err := svc.GetQueueUrl(&sqs.GetQueueUrlInput{
		QueueName: &queueName,
	})

	if err != nil {
		errMessage := fmt.Sprintf("Get queue url error: %s", err.Error())
		fmt.Println(errMessage)
		return nil, errors.New(errMessage)
	}

	queueUrl := queueUrlOutput.QueueUrl
	if queueUrl == nil {
		errMessage := "Get queue url error: queue url nil"
		fmt.Println(errMessage)
		return nil, errors.New(errMessage)
	}

	// set các giá trị tối đa / tối thiểu mà aws cho phép
	if delaySeconds < 0 {
		delaySeconds = 0
	} else if delaySeconds > 900 {
		delaySeconds = 900
	}

	if waitTime < 0 {
		waitTime = 0
	} else if waitTime > 20 {
		waitTime = 20
	}

	if maxNumberOfMsg <= 0 {
		maxNumberOfMsg = 1
	} else if maxNumberOfMsg > 10 {
		maxNumberOfMsg = 10
	}

	return &SqsWrapper{svc, *queueUrl, delaySeconds, maxNumberOfMsg, waitTime, visibilityTimeout}, nil
}

func (sqsWrapper SqsWrapper) SendStandardMsg(msg string) (*sqs.SendMessageOutput, error) {
	if sqsWrapper.Sqs == nil {
		return nil, nilSqsError
	}

	result, err := sqsWrapper.Sqs.SendMessage(&sqs.SendMessageInput{
		DelaySeconds: &sqsWrapper.DelaySeconds,
		MessageBody:  aws.String(msg),
		QueueUrl:     &sqsWrapper.QueueUrl,
	})

	if err != nil {
		return nil, err
	}

	return result, nil
}

func (sqsWrapper SqsWrapper) PullMessages() ([]*sqs.Message, error) {
	if sqsWrapper.Sqs == nil {
		return nil, nilSqsError
	}

	results, err := svc.ReceiveMessage(&sqs.ReceiveMessageInput{
		QueueUrl: &sqsWrapper.QueueUrl,
		MessageAttributeNames: aws.StringSlice([]string{
			"All",
		}),
		MaxNumberOfMessages: &sqsWrapper.MaxNumberOfMsg,
		WaitTimeSeconds:     &sqsWrapper.WaitTime,
		VisibilityTimeout:   &sqsWrapper.VisibilityTimeout,
	})
	// snippet-end:[sqs.go.send_receive_long_polling.call2]
	if err != nil {
		return nil, err
	}

	return results.Messages, nil
}

func (sqsWrapper SqsWrapper) DeleteMessage(message *sqs.Message) (*sqs.DeleteMessageOutput, error) {
	if sqsWrapper.Sqs == nil {
		return nil, nilSqsError
	}

	result, err := svc.DeleteMessage(&sqs.DeleteMessageInput{
		QueueUrl:      &sqsWrapper.QueueUrl,
		ReceiptHandle: message.ReceiptHandle,
	})

	if err != nil {
		return nil, err
	}

	return result, nil
}

func GetSqsMessageBody(message *sqs.Message) (*string, error) {
	if message.Body == nil {
		return nil, errors.New("Message's body nil")
	}

	return message.Body, nil
}
