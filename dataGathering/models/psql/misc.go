package psql

// SQL syntax helper structs

type tableMapper struct {
	columnNames string
	entity      interface{}
}

var TablesMapping = map[string]tableMapper{
	"sources": {
		columnNames: "(source_id, name)",
		entity:      &Source{},
	},
	"categories": {
		columnNames: "(category_id, name)",
		entity:      &Category{},
	},
	"podcasts": {
		columnNames: "(podcast_id, rss_link, last_update, title)",
		entity:      &Podcast{},
	},
	"itunes_ID": {
		columnNames: "(podcast_id, itunes_id, last_update)",
		entity:      &ItunesID{},
	},
}

type Request struct {
	Field interface{}
	Value interface{}
}
