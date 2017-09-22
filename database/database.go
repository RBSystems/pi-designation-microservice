package database

import (
	"database/sql"
	"log"
	"sync"

	"github.com/fatih/color"
)

/** lock things down here **/
var once sync.Once

/** all the good stuff lives here **/
var db *sql.DB

func DB() *sql.DB {
	once.Do(func() {
		//build source data
		data := "root:@tcp(localhost:3306)/"
		var err error

		db, err = sql.Open("mysql", data)
		if err != nil {
			log.Panicf("%s", color.HiRedString("[dbo] unable to open database: %s", err.Error()))
		}
	})

	return db
}
