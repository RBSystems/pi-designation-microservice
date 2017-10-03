package database

import (
	"database/sql"
	"log"
	"sync"

	_ "github.com/go-sql-driver/mysql" // Blank import due to its use as a driver

	"github.com/fatih/color"
)

/** lock things down here **/
var once sync.Once

/** all the good stuff lives here **/
var db *sql.DB

func DB() *sql.DB {
	once.Do(func() {
		//build source data
		data := "root:cmonster@tcp(localhost:3306)/room_designation"
		var err error

		db, err = sql.Open("mysql", data)
		if err != nil {
			log.Panicf("%s", color.HiRedString("[dbo] unable to open database: %s", err.Error()))
		}
	})

	return db
}
