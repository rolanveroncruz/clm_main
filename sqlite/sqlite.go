package sqlite

import (
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
	auth "ph.certs.com/clm_main/auth/sql"
	certs "ph.certs.com/clm_main/certs/sql"
)

func init() {
	InitializeDatabase()
	defer CloseDatabase()
}

var TheDB *sql.DB
var AuthQueryCental *auth.Queries
var CertsQueryCental *certs.Queries

func InitializeDatabase() {
	var err error
	TheDB, err = sql.Open("sqlite3", "/Volumes/RVC/Projects/certs.com.ph/clm_main/sqlite/auth")
	if err != nil {
		panic(err)
	}
	AuthQueryCental = auth.New(TheDB)
	CertsQueryCental = certs.New(TheDB)
}

func CloseDatabase() {
	err := TheDB.Close()
	if err != nil {
		return
	}
}
