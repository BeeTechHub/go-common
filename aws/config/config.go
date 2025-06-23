package awsConfig

import (
	"fmt"

	"github.com/aws/aws-sdk-go/aws/session"
)

var awsSession *session.Session

func InitAws() {
	awsSession = session.Must(session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
	}))
	//logger.Infof("init session aws:%s", awsSession)
}

func GetAWSSession() *session.Session {
	if awsSession == nil {
		InitAws()

		if awsSession == nil {
			fmt.Println("AWS Session null.")
		}
	}

	return awsSession
}
