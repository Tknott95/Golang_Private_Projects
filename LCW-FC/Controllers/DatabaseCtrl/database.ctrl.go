package databaseCtrl

import (
	"database/sql"

	_ "github.com/go-sql-driver/mysql"
	dbGlobals "github.com/tknott95/LCW-FC/Globals/db"
	sqlModel "github.com/tknott95/LCW-FC/Models/db"
)

var Store = newDB()

func newDB() *sqlModel.SQLStore {
	db, err := sql.Open("mysql", dbGlobals.SqlURL)
	if err != nil {
		println("🔒 Connection to AWS database established.\n")
	}

	return &sqlModel.SQLStore{
		DB: db,
	}
}
