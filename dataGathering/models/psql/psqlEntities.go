package psql

import (
	"podcastAnalyzer/parser/misc"
	"time"
)

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


var PostgresTablesMapping = map[string]misc.TableMapper{
	"sources": {
		ColumnNames: "(source_id, name)",
		Entity:      &Source{},
	},
	"categories": {
		ColumnNames: "(category_id, name)",
		Entity:      &Category{},
	},
	"podcasts": {
		ColumnNames: "(podcast_id, rss_link, last_update, title)",
		Entity:      &Podcast{},
	},
	"itunes_ID": {
		ColumnNames: "(podcast_id, itunes_id, last_update)",
		Entity:      &ItunesID{},
	},
}