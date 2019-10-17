package scraper_test

import (
	"strings"
	"testing"

	"github.com/TakumiOgawa/viewAnimeList/scraper"
)

func TestGetAnimeList(t *testing.T) {
	expect := "program"
	animeList := scraper.GetAnimeLink("2019", "spring")
	for _, anime := range animeList {
		if strings.Contains(anime, expect) {
			t.Errorf(`expect="%s" actual="%s"`, expect, anime)
		}
	}
}
