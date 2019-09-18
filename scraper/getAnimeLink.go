package scraper

import (
	"fmt"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/TakumiOgawa/viewAnimeList/util"
)

const (
	seasonURL string = "https://anime.eiga.com/program/season/"
)

func GetAnimeLink(year, season string) (contentURLArray []string) {

	reqURL := seasonURL + "/" + year + "-" + season

	doc, err := goquery.NewDocument(reqURL)
	if err != nil {
		fmt.Print("url scarapping failed")
	}
	// アニメコンテンツの部分のみのdivを取得
	selection := doc.Find("#mainContentsWide > div:nth-child(1) > div.articleInner > div > div ").Children()
	fmt.Println(selection.Length())

	selection.Find("a").Each(func(_ int, s *goquery.Selection) {
		url, _ := s.Attr("href")
		if strings.Contains(url, "program") {
			contentURLArray = append(contentURLArray, "https://anime.eiga.com"+url)
		}

	})

	contentURLArray = util.RemoveDupInSlice(contentURLArray)

	return
}
