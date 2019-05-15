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
	thread2Count := flag.Int("thread2", 0, "Wether we use parallelism or not")
	stupidCount := flag.Int("stupid", 0, "Wether we use parallelism or not")
	loopFlag := flag.Int("loop", 1, "Number of loops")
	flag.Parse()

	config := readConf()

	db := dbConnect(config)
	defer db.Close()

	for _, tableName := range config.Queries.Count {
		for i := 0; i < *loopFlag; i++ {
			if *threadCount > 0 {
				parallelCount(db, tableName, *threadCount)
			} else if *thread2Count > 0 {
				parallelCountV2(db, tableName, *thread2Count)
			} else if *stupidCount > 0 {
				stupidCounter(db, tableName)
			} else {
				coolCounter(db, tableName)
			}
		}
	}
}
