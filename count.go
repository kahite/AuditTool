package main

import (
	"database/sql"
	"fmt"
	"strings"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

func getDBNames(db *sql.DB, tableName string) []string {
	// t0 := time.Now()
	var dbNames []string

	query := fmt.Sprintf("SELECT table_schema FROM information_schema.tables WHERE table_name = \"%s\"", tableName)

	rows, err := db.Query(query)
	if err != nil {
		panic(err.Error()) // proper error handling instead of panic in your app
	}
	defer rows.Close()

	for rows.Next() {
		database := ""
		if err := rows.Scan(&database); err == nil {
			dbNames = append(dbNames, database)
		}
	}

	// t1 := time.Now()
	// fmt.Printf("DB names: %v\n", t1.Sub(t0))

	return dbNames
}

func stupidCounter(db *sql.DB, tableName string) {
	t0 := time.Now()

	dbNames := getDBNames(db, tableName)
	queryCount := 0

	for _, name := range dbNames {
		query := fmt.Sprintf("SELECT count(*) FROM %s.%s", name, tableName)
		queryRows, err := db.Query(query)

		if err != nil {
			continue
		}

		for queryRows.Next() {
			count := 0
			queryRows.Scan(&count)
			queryCount += count
		}
	}

	t1 := time.Now()
	fmt.Printf("Stupid counter %s: %d results in %v\n", tableName, queryCount, t1.Sub(t0))
}

func coolCounter(db *sql.DB, tableName string) {
	t0 := time.Now()
	dbNames := getDBNames(db, tableName)
	var dbQueries []string
	query := ""

	for _, name := range dbNames {
		dbQueries = append(dbQueries, fmt.Sprintf("SELECT count(*) AS countValue FROM `%s`.`%s`", name, tableName))
		query = strings.Join(dbQueries, " UNION ALL ")
	}
	query = "SELECT SUM(countValue) FROM (" + query + ") coolQuery"

	queryRows, _ := db.Query(query)
	queryCount := 0

	for queryRows.Next() {
		count := 0
		queryRows.Scan(&count)
		queryCount += count
	}

	t1 := time.Now()
	fmt.Printf("Cool counter %s: %d results in %v\n", tableName, queryCount, t1.Sub(t0))
}

func parallelCount(db *sql.DB, tableName string, threadCount int) {
	t0 := time.Now()

	dbNames := getDBNames(db, tableName)
	dbCount := len(dbNames)
	queryCount := 0
	coreCount := threadCount //runtime.NumCPU()
	coreChan := make(chan int, coreCount)

	for i := 0; i < coreCount; i++ {
		i := i
		go func() {
			innerCount := 0
			for j := i * dbCount / coreCount; j < (i+1)*dbCount/coreCount; j++ {
				query := fmt.Sprintf("SELECT count(*) FROM %s.%s", dbNames[j], tableName)
				queryRows, err := db.Query(query)
				if err != nil {
					continue
				}
				for queryRows.Next() {
					count := 0
					queryRows.Scan(&count)
					innerCount += count
				}
			}

			coreChan <- innerCount
		}()
	}

	for i := 0; i < coreCount; i++ {
		queryCount += <-coreChan
	}

	t1 := time.Now()
	fmt.Printf("Parallel counter %s: %d results in %v\n", tableName, queryCount, t1.Sub(t0))
}
