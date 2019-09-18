package main

import (
	"encoding/json"
	"net/http"
	"os"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/nlopes/slack/slackevents"
)

func main() {
	lambda.Start(handler)
}

func handler(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {

	eventsAPIEvent, _ := slackevents.ParseEvent(json.RawMessage(request.Body),
		slackevents.OptionVerifyToken(&slackevents.TokenComparator{VerificationToken: os.Getenv("V_TOKEN")}))

	if eventsAPIEvent.Type == slackevents.URLVerification {
		var r *slackevents.ChallengeResponse
		err := json.Unmarshal([]byte(request.Body), &r)
		if err != nil {
			return events.APIGatewayProxyResponse{Body: request.Body, StatusCode: http.StatusInternalServerError}, nil
		}
		return events.APIGatewayProxyResponse{
			Body: r.Challenge,
			Headers: map[string]string{
				"Content-Type": "text"},
			StatusCode: 200,
		}, nil
	}
	return events.APIGatewayProxyResponse{
		Body: "失敗",
		Headers: map[string]string{
			"Content-Type": "text"},
		StatusCode: 400,
	}, nil
}
