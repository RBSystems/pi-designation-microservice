package database

import (
	"log"
	"os"
	"sync"

	"github.com/fatih/color"
	_ "github.com/go-sql-driver/mysql" // Blank import due to its use as a driver
	"github.com/jmoiron/sqlx"
)

/** lock things down here **/
var once sync.Once

/** all the good stuff lives here **/
var db *sqlx.DB

func DB() *sqlx.DB {
	once.Do(func() {
		//build source data
		data := os.Getenv("DESIGNATION_DATABASE_USERNAME") + ":" +
			os.Getenv("DESIGNATION_DATABASE_PASSWORD") + "@tcp(" +
			os.Getenv("DESIGNATION_DATABASE_HOST") + ":" +
			os.Getenv("DESIGNATION_DATABASE_PORT") + ")" + "/" +
			os.Getenv("DESIGNATION_DATABASE_NAME")

		log.Printf("%s", color.HiCyanString("[database] data: %s", data))
		db = sqlx.MustOpen("mysql", data)
	})

	return db
}
