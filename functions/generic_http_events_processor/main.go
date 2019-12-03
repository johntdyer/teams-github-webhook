package main

import (
	"bytes"
	"context"
	"os"

	"github.com/aws/aws-sdk-go/service/sqs"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sns"

	"encoding/json"

	log "github.com/sirupsen/logrus"
)

var (
	snsHandler  *sns.SNS
	svc         *sqs.SQS
	emptyHeader = map[string]string{}
)

func init() {
	log.SetFormatter(&log.JSONFormatter{})
	log.SetOutput(os.Stdout)
	log.SetLevel(log.DebugLevel)
}

func main() {
	lambda.Start(handleRequest)
}

func setupConnection() {
	log.Info("setup sqs connection")

	sess := session.Must(session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
	}))

	svc = sqs.New(sess)
	snsHandler = sns.New(sess)
	log.Debug("Connection setup")

}

func apiGatewayResponse(str string, code int, headers map[string]string, body string) (events.APIGatewayProxyResponse, error) {
	return events.APIGatewayProxyResponse{ // Error HTTP response
		Body:       str,
		StatusCode: code,
	}, nil
}

// Handler is our lambda handler invoked by the `lambda.Start` function call
func handleRequest(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	log.Debug("### Event received ")
	var buf bytes.Buffer

	if !checkMAC(request.Body, request.Headers["X-Spark-Signature"], os.Getenv("HTTP_WEBHOOK_SECRET")) {
		log.Warnf("WARNING - MESSAGE RECEIVED WITH INVALID AUTH HEADER: " + request.Headers["X-Spark-Signature"])
		return apiGatewayResponse("Invalid message signature", 400, emptyHeader, "")
	}

	b := bytes.NewBufferString(request.Body)

	decoder := json.NewDecoder(b)
	var t = &webhookData{}

	err := decoder.Decode(&t)
	if err != nil {
		return apiGatewayResponse(err.Error(), 400, emptyHeader, "")
	}

	attachment := &attachmentData{}
	// decode the data and return actionData
	err = attachment.decodeTeamsData(t.Data.ID)
	if err != nil {
		panic(err)
	}

	// output := &outputEvent{}
	err = attachment.SendSQSEvent()

	if err != nil {
		log.Errorf("#### ERRROR MAin: %s\n", err)
		return apiGatewayResponse(err.Error(), 502, emptyHeader, "")
	}

	body, _ := json.Marshal(map[string]interface{}{
		"message": "accepted",
	})
	json.HTMLEscape(&buf, body)

	return events.APIGatewayProxyResponse{
		StatusCode:      200,
		IsBase64Encoded: false,
		Body:            buf.String(),
		Headers: map[string]string{
			"Content-Type": "application/json",
		},
	}, nil

}
