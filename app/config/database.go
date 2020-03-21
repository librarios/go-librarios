package config

import (
	"errors"
	"fmt"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	"log"
)

var DB *gorm.DB

// InitDB initializes DB connection
func InitDB(props map[string]interface{}) error {
	dialect, ok := props["dialect"]
	if !ok {
		return errors.New("DB dialect is not set")
	}

	strDialect := dialect.(string)

	var db *gorm.DB
	var err error

	switch strDialect {
	case "sqlite3":
		filename := props["filename"].(string)
		db, err = connectSqlite3(filename)

	case "mariadb":
		fallthrough
	case "mysql":
		db, err = connectMysql(
			props["host"].(string),
			props["database"].(string),
			props["username"].(string),
			props["password"].(string),
			props["port"].(int),
		)

	default:
		return errors.New(fmt.Sprintf("unsupported dialect: %s", dialect))
	}

	if err != nil {
		return err
	}

	// log mode
	if props["showSql"] == true {
		db.LogMode(true)
	}

	// set global DB connection
	DB = db

	return nil
}

// CloseDB closes DB connection
func CloseDB() {
	if DB != nil {
		_ = DB.Close()
	}
}

// connectSqlite3 connects to sqlite3 database
func connectSqlite3(filename string) (*gorm.DB, error) {
	db, err := gorm.Open("sqlite3", filename)
	if err != nil {
		return nil, err
	}
	log.Printf("connected to sqlite3: %s", filename)

	return db, nil
}

// connectMysql connects to mysql/mariaDB database
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
