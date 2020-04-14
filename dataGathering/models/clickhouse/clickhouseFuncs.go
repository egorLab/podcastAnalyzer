package clickhouse

import (
	"fmt"
	_ "github.com/ClickHouse/clickhouse-go"
	"github.com/jmoiron/sqlx"
	"log"
	"podcastAnalyzer/parser"
	"podcastAnalyzer/parser/logging"
)

type InstanceClickhouse struct {
	DB *sqlx.DB // connection
}

func NewClickhouseConnection() InstanceClickhouse {
	db, err := sqlx.Open("clickhouse", "tcp://127.0.0.1:9000?username=pdcst&password=pdcst&database=stats")
	if err != nil {
		logging.CheckErr(err, "Can't connect to the database")
	}

	err = db.Ping()
	if err != nil {
		logging.CheckErr(err, "Can't ping the database")
	}
	instance := InstanceClickhouse{DB: db}
	InitClickhouseTables(&instance)
	fmt.Println("Successfully connected clickhouse")
	return instance
}

func (clickhouse *InstanceClickhouse) InsertIntoTable(tablename string, entry interface{}) {
	mapping := dataGathering.ClickhouseTablesMapping[tablename]
	fieldsToFill, wildcard := dataGathering.GetFieldsAndWildcards(entry)

	sqlStatement := "INSERT INTO " + tablename + " " + mapping.ColumnNames + " VALUES " + wildcard

	_, err := clickhouse.DB.Exec(sqlStatement, fieldsToFill...)

	if err != nil {
		logging.CheckErr(err, sqlStatement)
	} else {
		logging.Logger.Info("Insert sent: ", fieldsToFill)
	}

	defer clickhouse.DB.Close()
}

func InitClickhouseTables(clickhouse *InstanceClickhouse) {
	_, err := clickhouse.DB.Exec(`
CREATE TABLE IF NOT EXISTS Podcasts (
			podcast_id    UInt64,
			main_category UInt16,
			all_main_categories   Array(Int16),
			title   String,
			listens_count   UInt64,
			comments_count  UInt64,
			rating UInt16,
			episodes_count UInt16,
			timestamp 	DateTime,
			source UInt8
		) engine=Memory
`)
	_, err = clickhouse.DB.Exec(`
CREATE TABLE IF NOT EXISTS Episodes (
			podcast_id    UInt64,
			episode_id UInt16,
			title   String,
			description String,
			length UInt16,
			listens_count   UInt64,
			comments_count  UInt64,
			trending_words	Array(String),
			rating UInt16,	
			publication_date 	DateTime,
			timestamp 	DateTime,
			explicit UInt8,
			is_trailer UInt8,
			timecodes_count UInt16,
			parts_count UInt8,
			source UInt8
		) engine=Memory
`)
	if err != nil {
		log.Fatal(err)
	}
}
