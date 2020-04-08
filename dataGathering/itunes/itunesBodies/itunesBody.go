package itunesBodies

// from https://github.com/ko/feedparser/blob/ecfd72b7f65820490fb93d5efad154aea216ab0f/search.go

type ItunesSearchResultsItem struct {
	ArtistID               int      `json:"artistId"`
	ArtistName             string   `json:"artistName"`
	ArtistViewURL          string   `json:"artistViewUrl"`
	ArtworkURL100          string   `json:"artworkUrl100"`
	ArtworkURL30           string   `json:"artworkUrl30"`
	ArtworkURL60           string   `json:"artworkUrl60"`
	ArtworkURL600          string   `json:"artworkUrl600"`
	CollectionCensoredName string   `json:"collectionCensoredName"`
	CollectionExplicitness string   `json:"collectionExplicitness"`
	CollectionHdPrice      float32  `json:"collectionHdPrice"`
	CollectionID           int      `json:"collectionId"`
	CollectionName         string   `json:"collectionName"`
	CollectionPrice        float32  `json:"collectionPrice"`
	CollectionViewURL      string   `json:"collectionViewUrl"`
	ContentAdvisoryRating  string   `json:"contentAdvisoryRating"`
	Country                string   `json:"country"`
	Currency               string   `json:"currency"`
	FeedURL                string   `json:"feedUrl"`
	GenreIds               []string `json:"genreIds"`
	Genres                 []string `json:"genres"`
	Kind                   string   `json:"kind"`
	PrimaryGenreName       string   `json:"primaryGenreName"`
	ReleaseDate            string   `json:"releaseDate"`
	TrackCensoredName      string   `json:"trackCensoredName"`
	TrackCount             int      `json:"trackCount"`
	TrackExplicitness      string   `json:"trackExplicitness"`
	TrackHdPrice           float32  `json:"trackHdPrice"`
	TrackHdRentalPrice     float32  `json:"trackHdRentalPrice"`
	TrackID                int      `json:"trackId"`
	TrackName              string   `json:"trackName"`
	TrackPrice             float32  `json:"trackPrice"`
	TrackRentalPrice       float32  `json:"trackRentalPrice"`
	TrackViewURL           string   `json:"trackViewUrl"`
	WrapperType            string   `json:"wrapperType"`
}

type ItunesSearchResults struct {
	ResultCount int                       `json:"resultCount"`
	Results     []ItunesSearchResultsItem `json:"results"`
}
