package psql

import "time"

type Category struct {
	CategoryID int    `db:"category_id"`
	Name       string `db:"name"`
}

type Source struct {
	SourceId int
	Name     string
}

type ItunesID struct {
	PodcastID  int
	ItunesID   int
	LastUpdate time.Time
}

type Podcast struct {
	PodcastID  int
	RSSLink    string
	LastUpdate time.Time
	Title      string
}
