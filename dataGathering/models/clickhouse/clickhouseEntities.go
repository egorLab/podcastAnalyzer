package clickhouse

import (
	"podcastAnalyzer/parser/misc"
	"time"
)
import _ "podcastAnalyzer/parser/misc"

var ClickhouseTablesMapping = map[string]misc.TableMapper{
	"Podcasts": {
		ColumnNames: "(podcast_id, main_category, all_main_categories, title, " +
			"listens_count, comments_count, rating, episodes_count, timestamp, source)",
		Entity: &Podcast{},
	},
	"Episodes": {
		ColumnNames: "(podcast_id, episode_id, description, title, " +
			"length, listens_count, comments_count, trending_words, rating, " +
			"publication_date, timestamp, explicit, is_trailer, timecodes_count, parts_count, source)",
		Entity: &Episode{},
	},
}

type Podcast struct {
	PodcastID         uint64    `db:"podcast_id"`
	MainCategory      uint16    `db:"main_category"`
	AllMainCategories []int16   `db:"all_main_categories"`
	Title             string    `db:"title"`
	ListensCount      uint64    `db:"listens_count"`
	CommentsCount     uint64    `db:"comments_count"`
	Rating            uint16    `db:"rating"`
	EpisodesCount     uint16    `db:"episodes_count"`
	Timestamp         time.Time `db:"timestamp"`
	Source            uint8     `db:"source"`
}

type Episode struct {
	PodcastID       uint64    `db:"podcast_id"`
	EpisodeID       uint16    `db:"episode_id"`
	Description     string    `db:"description"`
	Title           string    `db:"title"`
	Length          uint16    `db:"length"`
	ListensCount    uint64    `db:"listens_count"`
	CommentsCount   uint64    `db:"comments_count"`
	TrendingWords   []string  `db:"trending_words"`
	Rating          uint16    `db:"rating"`
	PublicationDate time.Time `db:"publication_date"`
	Timestamp       time.Time `db:"timestamp"`
	Explicit        uint8     `db:"explicit"`
	IsTrailer       uint8     `db:"is_trailer"`
	TimecodesCount  uint64    `db:"timecodes_count"`
	PartsCount      uint64    `db:"parts_count"`
	Source          uint8     `db:"source"`
}


