package main

import (
	b64 "encoding/base64"
	"fmt"
	"math/rand"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sqs"
)

type whSession struct {
	svc *sqs.SQS
	url string
}

func setupReceiveSession(cfg relayConfig) whSession {
	sess := session.Must(session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
		Config: aws.Config{
			Region:      aws.String(cfg.awsRegion),
			Credentials: credentials.NewStaticCredentials(cfg.awsAccessId, cfg.awsSecret, ""),
		},
	}))

	svc := sqs.New(sess)

	qURL := cfg.awsUrl
	return whSession{svc: svc, url: qURL}
}

func receiveMsgOnSession(whs whSession, timeout int64) (*string, string) {
	result, err := whs.svc.ReceiveMessage(&sqs.ReceiveMessageInput{
		AttributeNames: []*string{
			aws.String(sqs.MessageSystemAttributeNameSentTimestamp),
		},
		MessageAttributeNames: []*string{
			aws.String(sqs.QueueAttributeNameAll),
		},
		QueueUrl:            &whs.url,
		MaxNumberOfMessages: aws.Int64(1),
		VisibilityTimeout:   aws.Int64(20), // 20 seconds
		WaitTimeSeconds:     aws.Int64(timeout),
	})

	if err != nil {
		fmt.Println("Error", err)
		return nil, ""
	}

	if len(result.Messages) == 0 {
		fmt.Println("Received no messages")
		return nil, ""
	}
	//fmt.Println("Result: ", result)

	resultDelete, err := whs.svc.DeleteMessage(&sqs.DeleteMessageInput{
		QueueUrl:      &whs.url,
		ReceiptHandle: result.Messages[0].ReceiptHandle,
	})

	if err != nil {
		fmt.Println("Delete Error", err)
		return nil, ""
	}
        enc, ok := result.Messages[0].MessageAttributes["Encoding"]
        encStr := ""
        if ok {
            encStr = enc.GoString()
        }

        _ = resultDelete
        _ = encStr
	//fmt.Println("Message Deleted", resultDelete, encStr)

	return result.Messages[0].Body, "base64"
}
