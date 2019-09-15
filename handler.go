package main

import (
	"github.com/TakumiOgawa/viewAnimeList/scraper"
	"github.com/TakumiOgawa/viewAnimeList/slack"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/cloudwatchlogs"

	"fmt"
	"os"
)

func main() {
	scraper.GetAnimeLink("2018", "summer")
	// lambda.Start(handler)
}

func handler() {
	sess := session.Must(session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
	}))

	svc := cloudwatchlogs.New(sess)

	// Get up to the last 100 log events for LOG-STREAM-NAME
	// in LOG-GROUP-NAME:
	resp, err := svc.GetLogEvents(&cloudwatchlogs.GetLogEventsInput{
		Limit:         aws.Int64(100),
		LogGroupName:  aws.String("/var/log/messages"),
		LogStreamName: aws.String("testStream"),
	})
	if err != nil {
		fmt.Println("Got error getting log events:")
		fmt.Println(err.Error())
		os.Exit(1)
	}

	gotToken := ""
	nextToken := ""

	for _, event := range resp.Events {
		gotToken = nextToken
		nextToken = *resp.NextForwardToken

		if gotToken == nextToken {
			break
		}

		slackUrl := os.Getenv("slackURL")
		username := "test"
		iconEmoji := ":nick:"
		iconURL := ""
		channel := "#random"
		slack := slack.NewSlack(slackUrl, *event.Message, username, iconEmoji, iconURL, channel)
		slack.Send()
	}
}
