package main

import (
	"database/sql"
	"flag"

	_ "github.com/go-sql-driver/mysql"
)

func dbConnect(config ConfigParameter) *sql.DB {
	dsn := config.Mizar.User + ":" + config.Mizar.Password + "@tcp(" + config.Mizar.Host + ")/"
	db, err := sql.Open("mysql", dsn)

	if err != nil {
		panic(err.Error()) // Just for example purpose. You should use proper error handling instead of panic
	}
	return db
}

func main() {
	threadCount := flag.Int("thread", 0, "Wether we use parallelism or not")
	flag.Parse()

	config := readConf()

	db := dbConnect(config)
	defer db.Close()

	for _, tableName := range config.Queries.Count {
		if *threadCount > 0 {
			parallelCount(db, tableName, *threadCount)
		} else {
			coolCounter(db, tableName)
		}

		// stupidCounter(db, tableName)
	}
}
