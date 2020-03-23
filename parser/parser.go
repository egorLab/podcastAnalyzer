package main

import (
	"errors"
	"flag"
	"fmt"
	_ "github.com/mattn/go-sqlite3"
	"github.com/mmcdole/gofeed"
	"github.com/sirupsen/logrus"
	"net/http"
	"os"
	"podcastAnalyzer/parser/itunesPodcastSearch"
	"strings"
)

var log = logrus.New()
//var db *sql.DB
//var err error
//db, err := sql.Open("sqlite3", "./urlstorage")
//checkErr(err)

func handler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		switch r.URL.Path {
		case "/podcast-by-id":
			getPodcastById(w, r)
		case "/top-k-podcasts":
			topKPodcastsSearch(w, r)
		case "/podcast-by-query":
			getPodcastByQuery(w, r)
		}
	} else {
		_, _ = fmt.Fprintf(w, "Sorry, only GET is supported.")
	}

}

func getPodcastByQuery(w http.ResponseWriter, r *http.Request) {
	// just search itunes by query and find podcasts
	fp := gofeed.NewParser()
	queryRequested := r.URL.Query()["query"][0]
	results := itunesPodcastSearch.Search(queryRequested)
	if len(results) == 0 {
		checkErr(errors.New("no results found"), "no results found by query: " + queryRequested)
	}

	extractedFeed, err := fp.ParseURL(results[0].FeedURL)
	checkErr(err, "no feed extracted")
	fmt.Println(extractedFeed.Categories, extractedFeed.Description)

}

func topKPodcastsSearch(w http.ResponseWriter, r *http.Request) {
	// parse N topPodcasts and get top podcast ID feed
	KRequested := r.URL.Query()["k"][0]
	fp := gofeed.NewParser() // should it be more global?
	topPodcastfeed, err := fp.ParseURL("https://rss.itunes.apple.com/api/v1/us/podcasts/top-podcasts/all/" + KRequested + "/explicit.rss")
	checkErr(err, "no feed extracted. topPodcastfeed parsing failed")

	// get top podcast feed link
	splittedLink := strings.SplitAfter(topPodcastfeed.Items[0].Link, "/")
	podcastID := splittedLink[len(splittedLink)-1][2:]
	fmt.Println("top podcast id is ", podcastID)
}

func getPodcastById(w http.ResponseWriter, r *http.Request) {
	IdRequested := r.URL.Query()["id"][0]
	results := itunesPodcastSearch.Search(IdRequested)

	fp := gofeed.NewParser()

	if len(results) == 0 {
		checkErr(errors.New("no results found"), "no results found by podcastID: " + IdRequested)
	}

	extractedFeed, err := fp.ParseURL(results[0].FeedURL)
	fmt.Println(extractedFeed.Categories, extractedFeed.Description)
	checkErr(err, "no feed extracted")

	// insert into db
	// stmt, err := db.Prepare("INSERT INTO urlStorage(podcastId, url) values(?,?)")
	//checkErr(err)
	//
	//res, err := stmt.Exec(IdRequested, results[0].FeedURL)
	//checkErr(err)

}

func main() {
	// create logger
	log.Out = os.Stdout
	file, err := os.OpenFile("../logrus.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)

	if err == nil {
	 log.Out = file
	} else {
	 log.Info("Failed to log to file, using default stderr")
	}

	log.Info("Started logging...")

	var port string
	http.HandleFunc("/", handler)
	flag.StringVar(&port, "port", "8000", "port to run on. Default 8000.")
	log.Fatal(http.ListenAndServe("localhost:" + port, nil))
}

func checkErr(err error, details string) {
	if err != nil {
		log.WithFields(logrus.Fields{
			"err": err,
			"details": details,
		}).Fatal()
	}
}