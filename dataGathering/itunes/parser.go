package main

import (
	"errors"
	"fmt"
	_ "github.com/mattn/go-sqlite3"
	"github.com/mmcdole/gofeed"
	"podcastAnalyzer/parser/itunes/itunesPodcastSearch"
	"podcastAnalyzer/parser/logging"
	"podcastAnalyzer/parser/models/clickhouse"
	"strings"
	"time"
)

func getPodcastByQuery(queryRequested string) {
	// just search itunes by query and find podcasts
	fp := gofeed.NewParser()
	results := itunesPodcastSearch.Search(queryRequested)
	if len(results) == 0 {
		logging.CheckErr(errors.New("no results found"), "no results found by query: "+queryRequested)
	}

	extractedFeed, err := fp.ParseURL(results[0].FeedURL)
	logging.CheckErr(err, "no feed extracted")
	fmt.Println(extractedFeed.Categories, extractedFeed.Description)

}

func topKPodcastsSearch(k int) {
	// parse N topPodcasts and get top podcast ID feed
	fp := gofeed.NewParser()
	topPodcastfeed, err := fp.ParseURL("https://rss.itunes.apple.com/api/v1/us/podcasts/top-podcasts/all/" + string(k) + "/explicit.rss")
	logging.CheckErr(err, "no feed extracted. topPodcastfeed parsing failed")

	// get top podcast feed link
	splittedLink := strings.SplitAfter(topPodcastfeed.Items[0].Link, "/") // TODO fix podcast ID extraction
	podcastID := splittedLink[len(splittedLink)-1][2:]
	fmt.Println("top podcast id is ", podcastID)
}

func getPodcastById(IdRequested string) *gofeed.Feed {
	// 360084272
	results := itunesPodcastSearch.Search(IdRequested)

	fp := gofeed.NewParser()

	if len(results) == 0 {
		logging.CheckErr(errors.New("no results found"), "no results found by podcastID: "+IdRequested)
	}

	extractedFeed, err := fp.ParseURL(results[0].FeedURL)
	fmt.Println(extractedFeed.Categories, extractedFeed.Description, extractedFeed.Updated)
	logging.CheckErr(err, "no feed extracted")
	return extractedFeed
}

func main() {
	logging.InitLogger(false)

	//c := psql.NewPsqlConnection()
	//extractedFeed := getPodcastById("360084272")
	//t, _ := time.Parse(time.RFC1123Z, "Tue, 07 Apr 2020 23:57:34 +0000")
	//p := psql.Podcast{
	//	PodcastID:  36008422,
	//	RSSLink:    extractedFeed.Link,
	//	LastUpdate: t,
	//	Title:      extractedFeed.Title,
	//}
	//c.InsertIntoTable("podcasts", p)
	p := clickhouse.Podcast{
		PodcastID:         0,
		MainCategory:      0,
		AllMainCategories: nil,
		Title:             "",
		ListensCount:      0,
		CommentsCount:     0,
		Rating:            0,
		EpisodesCount:     0,
		Timestamp:         time.Time{},
		Source:            0,
	}
	c := clickhouse.NewClickhouseConnection()
	c.InsertIntoTable("Podcasts", p)

}
