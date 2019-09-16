package main

import (
	"encoding/json"
	"errors"
	"log"
	"net/url"
	"strings"

	"github.com/TakumiOgawa/viewAnimeList/scraper"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

var (
	// ErrNon200Response non 200 status code in response
	ErrNon200Response = errors.New("Non 200 Response found")
)

type SlackMessage struct {
	Text string `json:"text"`
}

func main() {
	lambda.Start(handler)
}

func handler(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {

	var tw string
	var un string
	for _, value := range strings.Split(request.Body, "&") {
		param := strings.Split(value, "=")
		if param[0] == "trigger_word" {
			tw, _ = url.QueryUnescape(param[1])
		}
		if param[0] == "user_name" {
			un, _ = url.QueryUnescape(param[1])
		}
	}

	if un == "slackbot" {
		return events.APIGatewayProxyResponse{}, nil
	}

	var season string
	if tw == "春" || tw == "はる" {
		season = "spring"
	} else if tw == "夏" || tw == "なつ" {
		season = "summer"
	} else if tw == "秋" || tw == "あき" {
		season = "autumn"
	} else if tw == "冬" || tw == "ふゆ" {
		season = "winter"
	}

	animeList := scraper.GetAnimeLink("2018", season)
	var outputAnime string
	for _, anime := range animeList {
		outputAnime = outputAnime + anime + "\n"
	}

	j, err := json.Marshal(SlackMessage{Text: outputAnime})
	if err != nil {
		log.Print(err)
		return events.APIGatewayProxyResponse{Body: "エラー"}, err
	}

	return events.APIGatewayProxyResponse{
		Body:       string(j),
		StatusCode: 200,
	}, nil

}
