package sqlite

import (
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
	auth "ph.certs.com/clm_main/auth/sql"
)

func init() {
	InitializeDatabase()
	defer CloseDatabase()
}

var TheDB *sql.DB
var QueryCental *auth.Queries

func InitializeDatabase() {
	var err error
	TheDB, err = sql.Open("sqlite3", "/Volumes/RVC/Projects/certs.com.ph/clm_main/sqlite/auth")
	if err != nil {
		panic(err)
	}
	QueryCental = auth.New(TheDB)
}

func CloseDatabase() {
	err := TheDB.Close()
	if err != nil {
		return
	}
}
