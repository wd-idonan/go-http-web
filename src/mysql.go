package main


import (
	"database/sql"
	"log"
	_ "github.com/go-sql-driver/mysql"
)

const (
	dbMaxOpenConns = 128
	dbMaxIdleConns = 32
)

func initDB(mysqlDSN string)(*sql.DB, error) {
	dsn := mysqlDSN
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		log.Printf("open db failed, err: %v", err)
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		log.Printf("ping db failed, err: %s", dsn)
		return nil, err
	}

	db.SetMaxOpenConns(dbMaxOpenConns)
	db.SetMaxIdleConns(dbMaxIdleConns)
	log.Print("connect to db success, " + dsn)
	return db, nil
}