package main

import (
	"database/sql"
	"fmt"
	"os"
	_ "reflect"

	_ "github.com/go-sql-driver/mysql"
	// _ "github.com/golang-migrate/migrate/v4"
	// I'd like to use this library https://github.com/golang-migrate/migrate/blob/master/database/postgres/TUTORIAL.md,
	// but failed
)

func InitDB(dbUser string, dbPass string, dbHost string, dbPort string, dbname string) *sql.DB {
	db1, err := sql.Open("mysql", fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", dbUser, dbPass, dbHost, dbPort, ""))
	if err != nil {
		panic(err.Error())
	}

	_, err = db1.Exec("CREATE DATABASE  IF NOT EXISTS " + dbname)
	if err != nil {
		panic(err)
	}
	db1.Close()
	return OpenSQL(dbUser, dbPass, dbHost, dbPort, dbname)
}

func ExecuteFile(db *sql.DB, fileName string) sql.Result {
	data, err := os.ReadFile(fileName)
	if err != nil {
		panic(err.Error())
	}

	stmt, err := db.Exec(string(data))
	if err != nil {
		panic(err.Error())
	}
	return stmt

}

func OpenSQL(dbUser string, dbPass string, dbHost string, dbPort string, dbname string) *sql.DB {
	db, err := sql.Open("mysql", fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", dbUser, dbPass, dbHost, dbPort, dbname))
	if err != nil {
		panic(err.Error())
	}
	return db
}
