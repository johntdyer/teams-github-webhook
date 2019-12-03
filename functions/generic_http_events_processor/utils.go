package main

import (
	"bytes"
	b64 "encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"time"

	log "github.com/sirupsen/logrus"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/sqs"
	"github.com/go-resty/resty/v2"
)

// func getRoom(actor string)
func (r *attachmentData) decodeTeamsData(t string) error {

	client := resty.New()

	client.
		SetRetryCount(2).
		SetRetryWaitTime(2 * time.Second).
		SetRetryMaxWaitTime(9 * time.Second).
		SetRetryAfter(func(client *resty.Client, resp *resty.Response) (time.Duration, error) {
			return 0, errors.New("quota exceeded")
		})

	actionURL := fmt.Sprintf("https://api.ciscospark.com/v1/attachment/actions/%s", t)

	resp, err := client.R().
		SetHeader("Content-Type", "application/json").
		SetAuthToken(os.Getenv("TEAMS_BEARER_TOKEN")).
		Get(actionURL)

	if err != nil {
		return err
	}
	log.Debugf("Decode response Code: %d", resp.StatusCode())

	if resp.StatusCode() == 200 {

		b := bytes.NewBufferString(resp.String())

		decoder := json.NewDecoder(b)

		err := decoder.Decode(r)
		if err != nil {
			return err
		}
	}

	return nil
}

func (r *webexTeamsUser) fetchUser(personID string) error {
	client := resty.New()

	client.
		SetRetryCount(2).
		SetRetryWaitTime(2 * time.Second).
		SetRetryMaxWaitTime(9 * time.Second).
		SetRetryAfter(func(client *resty.Client, resp *resty.Response) (time.Duration, error) {
			return 0, errors.New("quota exceeded")
		})

	personURL := fmt.Sprintf("https://api.ciscospark.com/v1/people/%s", personID)

	resp, err := client.R().
		SetHeader("Content-Type", "application/json").
		SetAuthToken(os.Getenv("TEAMS_BEARER_TOKEN")).
		Get(personURL)

	if err != nil {
		return err
	}
	log.Debugf("getUser response Code: %d", resp.StatusCode())

	if resp.StatusCode() == 200 {

		b := bytes.NewBufferString(resp.String())

		decoder := json.NewDecoder(b)

		err := decoder.Decode(r)
		if err != nil {
			return err
		}
	}

	return nil
}

func (r *webexTeamsSpace) fetchSpace(spaceID string) error {
	client := resty.New()

	client.
		SetRetryCount(2).
		SetRetryWaitTime(2 * time.Second).
		SetRetryMaxWaitTime(9 * time.Second).
		SetRetryAfter(func(client *resty.Client, resp *resty.Response) (time.Duration, error) {
			return 0, errors.New("quota exceeded")
		})

	personURL := fmt.Sprintf("https://api.ciscospark.com/v1/rooms/%s", spaceID)

	resp, err := client.R().
		SetHeader("Content-Type", "application/json").
		SetAuthToken(os.Getenv("TEAMS_BEARER_TOKEN")).
		Get(personURL)

	if err != nil {
		return err
	}
	log.Debugf("fetchSpace response Code: %d", resp.StatusCode())

	if resp.StatusCode() == 200 {

		b := bytes.NewBufferString(resp.String())

		decoder := json.NewDecoder(b)

		err := decoder.Decode(r)
		if err != nil {
			return err
		}
	}

	return nil
}

//SendSQSEvent -  send an sqs event
func (r *attachmentData) SendSQSEvent() error {

	setupConnection()

	person := &webexTeamsUser{}
	err := person.fetchUser(r.PersonID)
	if err != nil {
		panic(err)
	}

	space := &webexTeamsSpace{}
	err = space.fetchSpace(r.RoomID)
	if err != nil {
		panic(err)
	}

	output := &outputEvent{
		Method: "Post",
		Body:   "",
		Headers: map[string]string{
			"pullRequestId": r.Inputs.PullRequestID,
			"reviewStatus":  r.Inputs.ReviewStatus,
			"comments":      r.Inputs.Comments,
			"repoName":      r.Inputs.RepoName,
			"actorID":       r.PersonID,
			"timestamp":     r.Created.String(),
			"spaceID":       r.RoomID,
			"personName":    person.DisplayName,
			"personEmail":   person.Emails[0],
			"personAvatar":  person.Avatar,
			"spaceName":     space.Title,
		},
	}

	opsEvent, err := json.Marshal(output)
	if err != nil {
		return err
	}

	input := &sqs.SendMessageInput{
		MessageBody: aws.String(b64.StdEncoding.EncodeToString([]byte(string(opsEvent)))),
		MessageAttributes: map[string]*sqs.MessageAttributeValue{
			"Encoding": &sqs.MessageAttributeValue{
				DataType:    aws.String("String"),
				StringValue: aws.String("base64"),
			},
		},
		QueueUrl: aws.String(fmt.Sprintf("https://sqs.%s.amazonaws.com/%s/%s", os.Getenv("AWS_REGION"), os.Getenv("ACCOUNT_ID"), os.Getenv("TARGET_SQS_QUEUE"))),
	}
	if os.Getenv("SQS_FIFO_QUEUE_ENABLED") == "true" {
		log.Debugf("#### Enabled fifif queue\n")
		input.SetMessageGroupId("httpEventProxy")
	}

	log.Debugf("###### SQS Event %+v", input)
	res, err := svc.SendMessage(input)
	if err != nil {
		log.Errorf("### ERROR: %s", err)
		return err
	}

	log.Debugf("#### SendMessage() - %s", res)
	return nil

}
