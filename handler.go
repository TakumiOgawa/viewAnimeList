package main

import (
	"github.com/TakumiOgawa/viewAnimeList/scraper"
	"github.com/TakumiOgawa/viewAnimeList/slack"
	"github.com/aws/aws-lambda-go/lambda"

	"os"
)

func main() {
	lambda.Start(handler)
}

func handler() {

	var year, season string
	contentURLArray := scraper.GetAnimeLink(year, season)

	slackUrl := os.Getenv("slackURL")
	username := "test"
	iconEmoji := ":nick:"
	iconURL := ""
	channel := "#random"

	for _, animeContentsURL := range contentURLArray {
		slack := slack.NewSlack(slackUrl, "https://anime.eiga.com"+animeContentsURL, username, iconEmoji, iconURL, channel)
		slack.Send()
	}

}
