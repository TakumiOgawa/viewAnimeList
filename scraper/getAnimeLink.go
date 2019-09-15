package scraper

import (
	"fmt"

	"github.com/PuerkitoBio/goquery"
)

const (
	seasonURL string = "https://anime.eiga.com/program/season/"
)

func GetAnimeLink(year, season string) {

	reqURL := seasonURL + "/" + year + "-" + season

	doc, err := goquery.NewDocument(reqURL)
	if err != nil {
		fmt.Print("url scarapping failed")
	}
	selection := doc.Find("#mainContentsWide > div:nth-child(1) > div.articleInner > div > div").Children()
	for _, node := range selection.Nodes {
		fmt.Println(node.Data)
	}

	// selection.Find("a").Each(func(_ int, s *goquery.Selection) {
	// 	url, _ := s.Attr("href")
	// 	fmt.Println(url)
	// })
}
