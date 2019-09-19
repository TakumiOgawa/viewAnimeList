package main

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/TakumiOgawa/viewAnimeList/scraper"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/nlopes/slack"
	"github.com/nlopes/slack/slackevents"
)

func main() {
	lambda.Start(handler)
}

func handler(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {

	requestBody := json.RawMessage(request.Body)
	slackeventsOption := slackevents.OptionVerifyToken(&slackevents.TokenComparator{VerificationToken: os.Getenv("V_TOKEN")})
	eventsAPIEvent, err := slackevents.ParseEvent(requestBody, slackeventsOption)

	if err != nil {
		log.Printf("%s err \n ", err)
	}

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

	if eventsAPIEvent.Type == slackevents.CallbackEvent {
		api := slack.New(os.Getenv("BOT_USER_TOKEN"))
		var year string
		var season string

		// デフォルトは現在年
		year = strconv.Itoa(time.Now().Year())

		innerEvent := eventsAPIEvent.InnerEvent
		switch ev := innerEvent.Data.(type) {
		case *slackevents.AppMentionEvent: // Botユーザーへのメンションの場合
			text := ev.Text

			switch {
			case strings.Contains(text, "春"):
				season = "spring"
			case strings.Contains(text, "夏"):
				season = "summer"
			case strings.Contains(text, "秋"):
				season = "autonum"
			case strings.Contains(text, "冬"):
				season = "winter"
			}

			animeList := scraper.GetAnimeLink(year, season)
			for _, anime := range animeList {
				_, _, err := api.PostMessage(os.Getenv("TARGET_CHANNEL_ID"), slack.MsgOptionText(anime, false))
				if err != nil {
					log.Printf("ERROR : %s", err)

				}
			}
		case *slackevents.MessageEvent:
			text := ev.Text
			if !strings.Contains(text, "メンションつけろやオラ") {
				reply := "メンションつけろやオラ"
				_, _, err := api.PostMessage(os.Getenv("TARGET_CHANNEL_ID"), slack.MsgOptionText(reply, false))
				if err != nil {
					log.Printf("ERROR : %s", err)

				}
			}
		}
	}

	return events.APIGatewayProxyResponse{
		Body: "失敗",
		Headers: map[string]string{
			"Content-Type": "text"},
		StatusCode: 400,
	}, nil
}
