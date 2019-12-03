package main

import (
	"context"
	"encoding/base64"
	"errors"
	"os"
	"strings"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"

	log "github.com/sirupsen/logrus"
)

func init() {
	log.SetFormatter(&log.JSONFormatter{})
	log.SetOutput(os.Stdout)
	log.SetLevel(log.DebugLevel)

}

func isAuthorized(username string, password string) bool {
	if strings.Compare(username, string(os.Getenv("BASIC_AUTH_USERNAME"))) != 0 {
		log.Printf("Invalid username received: %s", username)
		return false
	}

	if strings.Compare(password, string(os.Getenv("BASIC_AUTH_PASSWORD"))) != 0 {
		log.Printf("Bad password")
		return false
	}

	return true
}

func handler(ctx context.Context, request events.APIGatewayCustomAuthorizerRequest) (events.APIGatewayCustomAuthorizerResponse, error) {
	token := request.AuthorizationToken
	auth := strings.SplitN(token, " ", 2)

	if auth[0] == "bearer" {
		var bearerToken string
		if len(auth) > 1 {
			bearerToken = auth[len(auth)-1]
		}
		if bearerToken == os.Getenv("HTTP_AUTH_TOKEN") {
			return generatePolicy("user", "Allow", request.MethodArn), nil
		}
	}

	if len(auth) == 2 || auth[0] == "Basic" {

		payload, _ := base64.StdEncoding.DecodeString(auth[1])
		pair := strings.SplitN(string(payload), ":", 2)

		log.WithFields(log.Fields{
			"auth": auth,
		}).Warn("Auth Detected")

		if len(pair) == 2 && isAuthorized(pair[0], pair[1]) {
			return generatePolicy("user", "Allow", request.MethodArn), nil
		}

	}

	return events.APIGatewayCustomAuthorizerResponse{}, errors.New("Unauthorized")
}

func main() {
	lambda.Start(handler)
}

func generatePolicy(principalID, effect, resource string) events.APIGatewayCustomAuthorizerResponse {
	authResponse := events.APIGatewayCustomAuthorizerResponse{PrincipalID: principalID}

	if effect != "" && resource != "" {
		authResponse.PolicyDocument = events.APIGatewayCustomAuthorizerPolicy{
			Version: "2012-10-17",
			Statement: []events.IAMPolicyStatement{
				{
					Action:   []string{"execute-api:Invoke"},
					Effect:   effect,
					Resource: []string{resource},
				},
			},
		}
	}
	return authResponse
}
