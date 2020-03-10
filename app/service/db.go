package service

import (
	"errors"
	"fmt"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	"github.com/librarios/go-librarios/app/model"
	"log"
)

var dbConn *gorm.DB

func InitDB(props map[string]interface{}) error {
	conn, err := ConnectDB(props)
	if err != nil {
		return err
	}
	dbConn = conn
	return nil
}

func ConnectDB(props map[string]interface{}) (*gorm.DB, error) {
	dialect, ok := props["dialect"]
	if !ok {
		return nil, errors.New("DB dialect is not set")
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
		return nil, errors.New(fmt.Sprintf("unsupported dialect: %s", dialect))
	}

	if err != nil {
		return db, err
	}

	// log mode
	if props["showSql"] == true {
		db.LogMode(true)
	}

	// auto-migrate
	if props["autoMigrate"] == true {
		autoMigrateDB(db)
	}

	return db, err
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

func autoMigrateDB(db *gorm.DB) {
	db.AutoMigrate(
		&model.Book{},
		&model.OwnedBook{},
	)
	log.Println("DB auto migration finished.")
}
