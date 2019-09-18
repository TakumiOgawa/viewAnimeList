package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/TakumiOgawa/viewAnimeList/scraper"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/nlopes/slack"
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

	if eventsAPIEvent.Type == slackevents.CallbackEvent {
		api := slack.New(os.Getenv("BOT_USER_TOKEN"))
		var season string
		innerEvent := eventsAPIEvent.InnerEvent
		switch ev := innerEvent.Data.(type) {
		case *slackevents.AppMentionEvent: // Botユーザーへのメンションの場合
			reply := "応答テキスト"
			api.PostMessage(ev.Channel, slack.MsgOptionText(reply, false))
		case *slackevents.MessageEvent:
			text := ev.Text
			if text == "春" || text == "夏" || text == "秋" || text == "冬" {

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

				animeList := scraper.GetAnimeLink("2018", season)
				var outputAnime string
				for _, anime := range animeList {
					outputAnime = outputAnime + anime + "\n"
				}

				channelID, timestamp, err := api.PostMessage(os.Getenv("TARGET_CHANNEL_ID"), slack.MsgOptionText(outputAnime, false))
				if err != nil {
					fmt.Printf("%s\n", err)

				}

				fmt.Printf("Message successfully sent to channel %s at %s", channelID, timestamp)
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
