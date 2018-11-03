package config

import (
	"database/sql"
	"fmt"
	"log"
	"sync"
	_ "github.com/lib/pq"

)

type DBConnection struct {
	DBConn *sql.DB
}

var instance *DBConnection
var once sync.Once

func Connect() *DBConnection {
	conf, err := Load()
	if err != nil {
		fmt.Println("conf loading failed:", err)
	}

	db, err := sql.Open(conf.Cockroach.Dialect, fmt.Sprintf("host=%s port=%s user=%s dbname=%s sslmode=disable", conf.Cockroach.Host, conf.Cockroach.Port, conf.Cockroach.User, conf.Cockroach.DbName))
	once.Do(func() {
		fmt.Printf("host=%s port=%s user=%s dbname=%s password=%s", conf.Cockroach.Host, conf.Cockroach.Port, conf.Cockroach.User, conf.Cockroach.DbName, conf.Cockroach.Password)
		if err != nil {
			log.Fatal("error connecting to the database", err)
		}
		instance = &DBConnection{
			DBConn: db,
		}
	})
	return instance
}
