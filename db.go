package main

import (
	"errors"
	"fmt"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	"log"
)

var dbConn *gorm.DB

func ConnectDB(props Map) (*gorm.DB, error) {
	dialect, exists := props["dialect"]
	if !exists {
		return nil, errors.New("DB dialect is not set")
	}

	strDialect := dialect.(string)

	switch strDialect {
	case "sqlite3":
		filename := props["filename"].(string)
		return connectSqlite3(filename)

	case "mariadb":
		fallthrough
	case "mysql":
		return connectMysql(
			props["host"].(string),
			props["database"].(string),
			props["username"].(string),
			props["password"].(string),
			props["port"].(int),
			)

	default:
		return nil, errors.New(fmt.Sprintf("unsupported dialect: %s", dialect))
	}
}

func connectSqlite3(filename string) (*gorm.DB, error) {
	db, err := gorm.Open("sqlite3", filename)
	if err != nil {
		return nil, err
	}
	log.Printf("connected to sqlite3: %s", filename)

	return db, nil
}

func connectMysql(host, database, username, password string, port int) (*gorm.DB, error) {
	args := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&serverTimezone=UTC&parseTime=True",
		username,
		password,
		host,
		port,
		database,
	)

	db, err := gorm.Open("mysql", args)
	if err != nil {
		return nil, err
	}
	db.DB().SetMaxIdleConns(10)
	db.DB().SetMaxOpenConns(100)

	log.Printf("connected to mysql: %s@%s:%d/%s", username, host, port, database)

	return db, nil
}
