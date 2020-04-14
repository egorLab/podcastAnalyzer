package misc

import (
	"podcastAnalyzer/parser/models/clickhouse"
	"podcastAnalyzer/parser/models/psql"
	"reflect"
	"strconv"
)

// SQL syntax helper structs

type tableMapper struct {
	ColumnNames string
	Entity      interface{}
}

var PostgresTablesMapping = map[string]tableMapper{
	"sources": {
		ColumnNames: "(source_id, name)",
		Entity:      &psql.Source{},
	},
	"categories": {
		ColumnNames: "(category_id, name)",
		Entity:      &psql.Category{},
	},
	"podcasts": {
		ColumnNames: "(podcast_id, rss_link, last_update, title)",
		Entity:      &psql.Podcast{},
	},
	"itunes_ID": {
		ColumnNames: "(podcast_id, itunes_id, last_update)",
		Entity:      &psql.ItunesID{},
	},
}

var ClickhouseTablesMapping = map[string]tableMapper{
	"podcasts": {
		ColumnNames: "(podcast_id, main_category, all_main_categories, title, " +
			"listens_count, comments_count, rating, episodes_count, timestamp, source)",
		Entity: &clickhouse.Podcast{},
	},
	"episodes": {
		ColumnNames: "(podcast_id, episode_id, description, title, " +
			"length, listens_count, comments_count, trending_words, rating, " +
			"publication_date, timestamp, explicit, is_trailer, timecodes_count, parts_count, source)",
		Entity: &clickhouse.Episode{},
	},
}

type Request struct {
	Field interface{}
	Value interface{}
}

func GetFieldsAndWildcards(entry interface{}) (fields []interface{}, wildcardStr string) {
	// create interface to pass and wildcards
	wildcard := "("
	s := reflect.ValueOf(entry)
	fieldsToFill := make([]interface{}, s.NumField())
	for i := 0; i < s.NumField(); i++ {
		fieldsToFill[i] = s.Field(i).Interface()
		wildcard += "$" + strconv.Itoa(i+1)
		if i != s.NumField()-1 {
			wildcard += ", " // TODO better way to construct statement?
		}
	}
	wildcard += ")"

	return fieldsToFill, wildcard
}
