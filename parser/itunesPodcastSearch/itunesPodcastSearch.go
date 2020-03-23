package itunesPodcastSearch

import (
	"encoding/json"
	"log"
	"net/http"
	"podcastAnalyzer/parser/itunesBodies"
	"strconv"
	"time"
)

var myClient = &http.Client{Timeout: 10 * time.Second}


func SearchUrlGenerator(item string) (string, bool) {
	var searchString string
	var isIdBasedSearch bool
	if _, err := strconv.Atoi(item); err == nil {
		// if it is an id-based search
		searchString = "lookup?id=" + item
	} else {
		// if it is term-based search
		searchString = "search?term=" + item
	}

	return "https://itunes.apple.com/" + searchString + "&entity=podcast", isIdBasedSearch
}

func Search(query string) []itunesBodies.ItunesSearchResultsItem {
	var results itunesBodies.ItunesSearchResults

	// get url and flag if the search is ID based
	url, isIdBasedSearch := SearchUrlGenerator(query)
	err := GetJsonResultsFromURL(url, &results)
	if err != nil {
		log.Fatal(err)
	}

	if !isIdBasedSearch {
		// filter podcasts if search was by query
		results = FilterOnlyPodcasts(results)
	}

	return results.Results
}


func FilterOnlyPodcasts(results itunesBodies.ItunesSearchResults) itunesBodies.ItunesSearchResults {
	// from https://github.com/ko/feedparser/blob/ecfd72b7f65820490fb93d5efad154aea216ab0f/search.go

	var podcasts itunesBodies.ItunesSearchResults
	var podcastCount int
	var podcastList []itunesBodies.ItunesSearchResultsItem

	podcastCount = 0

	for index := range results.Results {
		podcastCount = podcastCount + 1
		if results.Results[index].Kind == "podcast" {
			podcastList = append(podcastList, results.Results[index])
		}
	}

	podcasts.Results = podcastList
	podcasts.ResultCount = podcastCount

	return podcasts
}

func GetJsonResultsFromURL(url string, target interface{}) error {
	r, err := myClient.Get(url)
	if err != nil {
		return err
	}
	defer r.Body.Close()

	return json.NewDecoder(r.Body).Decode(target)
}

