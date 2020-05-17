package main

import (
	"bufio"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"podcastAnalyzer/parser/models/clickhouse"
	"strconv"

	// "github.com/gorilla/mux"
	_ "github.com/mattn/go-sqlite3"
	"github.com/mmcdole/gofeed"
	"log"
	"net/http"
	"os"
	"os/signal"
	"podcastAnalyzer/parser/itunes/itunesPodcastSearch"
	"podcastAnalyzer/parser/logging"
	"strings"
	"syscall"
	"time"
)


func handler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		switch r.URL.Path {
		//case "/podcast-by-id":
		//	getPodcastById(w, r)
		//case "/top-k-podcasts":
		//	topKPodcastsSearch(w, r)
		//case "/podcast-by-query":
		//	getPodcastByQuery(w, r)
		case "/health":
			healthHandler(w, r)
		case "/readiness":
			readinessHandler(w, r)
		}
	} else {
		_, _ = fmt.Fprintf(w, "Sorry, only GET is supported.")
	}

}

func healthHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
}

func readinessHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
}

func getPodcastByQuery(queryRequested string) {
	// just search itunes by query and find podcasts
	fp := gofeed.NewParser()
	// queryRequested := extractAttribute(r, "query")
	results := itunesPodcastSearch.Search(queryRequested)
	if len(results) == 0 {
		logging.CheckErr(errors.New("no results found"), "no results found by query: "+queryRequested)
	}

	_, err := fp.ParseURL(results[0].FeedURL)
	logging.CheckErr(err, "no feed extracted")

	// _, err = w.Write([]byte(extractedFeed.Items[0].Description))
	logging.CheckErr(err, "writing feed failed")

}

func topKPodcastsSearch(k int) {
	// parse N topPodcasts and get top podcast ID feed
	// KRequested := extractAttribute(r, "k")
	fp := gofeed.NewParser()
	topPodcastfeed, err := fp.ParseURL("https://rss.itunes.apple.com/api/v1/us/podcasts/top-podcasts/all/" + string(k) + "/explicit.rss")

	logging.CheckErr(err, "no feed extracted. topPodcastfeed parsing failed")

	// get top podcast feed link
	splittedLink := strings.SplitAfter(topPodcastfeed.Items[0].Link, "/") // TODO fix podcast ID extraction
	_ = splittedLink[len(splittedLink)-1][2:]

	// _, err = w.Write([]byte(podcastID))
	logging.CheckErr(err, "writing feed failed")
}

func getPodcastById(IdRequested int)  (*gofeed.Feed, error) {
	// 360084272
	// IdRequested := extractAttribute(r, "id")
	results := itunesPodcastSearch.Search(strconv.Itoa(IdRequested))

	fp := gofeed.NewParser()

	if len(results) == 0 {
		logging.CheckErr(errors.New("no results found"), "no results found by podcastID: "+string(IdRequested))
	}
	feed, err := fp.ParseURL(results[0].FeedURL) // TODO fix me
	logging.CheckErr(err, "no feed extracted")

	// _, err = w.Write([]byte(extractedFeed.Description))
	// logging.CheckErr(err, "writing feed failed")

	//parsedTime, _ := time.Parse(time.RFC1123Z, feed.Updated)
	//
	//p := psql.Podcast{
	//	PodcastID:  IdRequested,
	//	RSSLink:    feed.Link,
	//	LastUpdate: parsedTime,
	//	Title:      feed.Title,
	//}

	return feed, err
	// instancePsql.InsertIntoTable("podcasts", p)
}

func extractAttribute(r *http.Request, attr string) string {
	return r.URL.Query()[attr][0]
}

func main() {
	logging.InitLogger(false)

	r := http.NewServeMux()

	// r.HandleFunc("/", handler)
	r.HandleFunc("/health", healthHandler)
	r.HandleFunc("/readiness", readinessHandler)

	srv := &http.Server{
		Handler:      r,
		Addr:         ":8080",
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}
	// c := psql.NewPsqlConnection()
	// getPodcastById(1480311435, &c)
	file, err := os.Open("podcast_ids.txt")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	category2id := make(map[string]int16)
	source2id := make(map[string]uint8)
	catID := 0
	authorID := 0
	c := clickhouse.NewClickhouseConnection()
	for scanner.Scan() {
		var currentPodcastCats []int16
		var mainCat uint16
		i, err := strconv.ParseInt(scanner.Text(), 10, 32)

		logging.CheckErr(err, "Can't convert")

		feed, err := getPodcastById(int(i))
		if err != nil {
			continue
		}

		parsedTime, _ := time.Parse(time.RFC1123, feed.Updated)

		for _, cat := range feed.Categories {
			if  _, ok := category2id[cat]; !ok  {
				category2id[cat] = int16(catID)
				catID += 1
			}
			currentPodcastCats = append(currentPodcastCats, category2id[cat])
		}
		if  _, ok := source2id[feed.Author.Name]; !ok  {
			source2id[feed.Author.Name] = uint8(authorID)
			authorID += 1
		}
		if len(currentPodcastCats) == 0 {
			mainCat = 9999
		} else {
			mainCat = uint16(currentPodcastCats[0])
		}

		p := clickhouse.Podcast{
			PodcastID:         uint64(i),
			MainCategory:      mainCat,
			AllMainCategories: currentPodcastCats,
			Title:             feed.Title,
			ListensCount:      0,
			CommentsCount:     0,
			Rating:            0,
			EpisodesCount:     uint16(len(feed.Items)),
			Timestamp:         parsedTime,
			Source:            source2id[feed.Author.Name],
		}

		c.InsertIntoTable("Podcasts", p)
	}

	jsonData, err := json.Marshal(category2id)

	if err != nil {
		panic(err)
	}

	jsonFile, err := os.Create("./category2id.json")

	if err != nil {
		panic(err)
	}
	defer jsonFile.Close()

	jsonFile.Write(jsonData)
	jsonFile.Close()

	jsonData, err = json.Marshal(source2id)

	if err != nil {
		panic(err)
	}

	jsonFile, err = os.Create("./source2id.json")

	if err != nil {
		panic(err)
	}
	defer jsonFile.Close()

	jsonFile.Write(jsonData)
	jsonFile.Close()


	go func() {
		logging.Logger.Info("Starting Server")
		if err := srv.ListenAndServe(); err != nil {
			logging.Logger.Fatal("Failed to start a server")
		}
	}()

	waitForShutdown(srv)
}


func waitForShutdown(srv *http.Server) {
	interruptChan := make(chan os.Signal, 1)
	signal.Notify(interruptChan, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	// Block until we receive our signal.
	<-interruptChan

	// Create a deadline to wait for.
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()
	err := srv.Shutdown(ctx)
	logging.CheckErr(err, "srv.Shutdown failed")

	log.Println("Shutting down")
	os.Exit(0)
}



//c := psql.NewPsqlConnection()
// extractedFeed := getPodcastById("360084272")
//t, _ := time.Parse(time.RFC1123Z, "Tue, 07 Apr 2020 23:57:34 +0000")
// fmt.Println(extractedFeed.Items[0])

//p := psql.Podcast{
//	PodcastID:  36008422,
//	RSSLink:    extractedFeed.Link,
//	LastUpdate: t,
//	Title:      extractedFeed.Title,
//}
//c.InsertIntoTable("podcasts", p)
//p := clickhouse.Podcast{
//	PodcastID:         0,
//	MainCategory:      0,
//	AllMainCategories: nil,
//	Title:             "",
//	ListensCount:      0,
//	CommentsCount:     0,
//	Rating:            0,
//	EpisodesCount:     0,
//	Timestamp:         time.Time{},
//	Source:            0,
//}
//c := clickhouse.NewClickhouseConnection()
//c.InsertIntoTable("Podcasts", p)



// CRON
//s1 := gocron.NewScheduler(time.UTC)
//
//s1.Every(3).Seconds().Do(getPodcastById, "360084272")
//
//// scheduler starts running jobs and current thread continues to execute
//<- s1.Start()