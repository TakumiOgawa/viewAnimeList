package main

import (
	"encoding/json"
	"errors"
	"log"
	"os"
	"strings"

	"github.com/TakumiOgawa/viewAnimeList/scraper"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/nlopes/slack"
	"github.com/nlopes/slack/slackevents"
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

	log.Println("method:", request.HTTPMethod, "path:", request.Path, "res:", request.Resource)
	eventsAPIEvent, _ := slackevents.ParseEvent(json.RawMessage(request.Body),
		slackevents.OptionVerifyToken(&slackevents.TokenComparator{VerificationToken: "TOKEN"}))

	if eventsAPIEvent.Type == slackevents.URLVerification {
		var r *slackevents.ChallengeResponse

		return events.APIGatewayProxyResponse{
			Body: r.Challenge,
			Headers: map[string]string{
				"Content-Type": "text"},
			StatusCode: 200,
		}, nil

		// こっからレスポンス
		api := slack.New(os.Getenv("SLACK_TOKEN"))
		var season string

		if eventsAPIEvent.Type == slackevents.CallbackEvent {
			innerEvent := eventsAPIEvent.InnerEvent
			switch ev := innerEvent.Data.(type) {
			case *slackevents.AppMentionEvent: // Botユーザーへのメンションの場合
				reply := "応答テキスト"
				api.PostMessage(ev.Channel, slack.MsgOptionText(reply, false))
			case *slackevents.MessageEvent:
				if ev.ChannelType == "im" { // ダイレクトメッセージの場合
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

					animeList := scraper.GetAnimeLink("2018", season)
					var outputAnime string
					for _, anime := range animeList {
						outputAnime = outputAnime + anime + "\n"
					}

					api.PostMessage("TARGET_CHANNEL_ID", slack.MsgOptionText(outputAnime, false))

				}
			}
		}
	}
	return events.APIGatewayProxyResponse{
		Body:       "badrequest",
		StatusCode: 400,
	}, nil
}
