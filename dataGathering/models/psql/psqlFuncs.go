package psql

import (
	"fmt"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"podcastAnalyzer/parser"
	"podcastAnalyzer/parser/logging"
)

const ( // TODO: put in config
	host     = "localhost"
	port     = 5432
	user     = "dev"
	password = "fsdf184-"
	dbname   = "postgres"
)

type InstancePsql struct {
	DB *sqlx.DB // connection
}

func NewPsqlConnection() InstancePsql {
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)
	db, err := sqlx.Open("postgres", psqlInfo)
	if err != nil {
		logging.CheckErr(err, "Can't connect to the database")
	}

	err = db.Ping()
	if err != nil {
		logging.CheckErr(err, "Can't ping the database")
	}
	instance := InstancePsql{DB: db}
	fmt.Println("Successfully connected postgres")
	return instance
}

func (psql *InstancePsql) InsertIntoTable(tablename string, entry interface{}) {
	mapping := dataGathering.PostgresTablesMapping[tablename]
	fieldsToFill, wildcard := dataGathering.GetFieldsAndWildcards(entry)

	sqlStatement := "INSERT INTO " + tablename + " " + mapping.ColumnNames + " VALUES " + wildcard

	_, err := psql.DB.Exec(sqlStatement, fieldsToFill...)

	if err != nil {
		logging.CheckErr(err, sqlStatement)
	} else {
		logging.Logger.Info("Insert sent: ", fieldsToFill)
	}

	defer psql.DB.Close()
}

func (psql *InstancePsql) GetRowFromTableWithWhere(tablename string, entry dataGathering.Request) interface{} {
	mapping := dataGathering.PostgresTablesMapping[tablename]
	sqlStatement := "SELECT * FROM " + tablename +
		" WHERE " + fmt.Sprintf("%v", entry.Field) + " = " + fmt.Sprintf("%v", entry.Value)

	row := psql.DB.QueryRowx("SELECT * FROM "+tablename+
		" WHERE "+fmt.Sprintf("%v ", entry.Field)+"= $1 ",
		entry.Value)

	entity := mapping.Entity
	err := row.StructScan(entity)
	if err != nil {
		logging.CheckErr(err, sqlStatement)
	} else {
		logging.Logger.Info("Query: ", sqlStatement)
	}

	defer psql.DB.Close()

	return entity
}
