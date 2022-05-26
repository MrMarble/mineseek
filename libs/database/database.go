package database

import (
	"database/sql"
	"log"
	"sync"

	_ "github.com/lib/pq" // postgre driver
	"github.com/mrmarble/mineseek/libs/minecraft"
)

var (
	data *DB
	once sync.Once
)

type DB struct {
	db *sql.DB
}

func (db *DB) InsertSLP(slp *minecraft.ServerListPing) error {
	_, err := db.db.Exec("INSERT INTO servers VALUES (%s,%d,%s,%s,%s,%d,%s)", slp.Host, slp.Port, slp.MOTD, slp.Favicon, slp.MOTD, slp.MaxPlayers, slp.MaxPlayers)
	return err
}

func getConnection() (*sql.DB, error) {
	return sql.Open("postgres", "") // TODO: Use env

}

func initDB() {
	db, err := getConnection()
	if err != nil {
		log.Panic(err)
	}
	data = &DB{
		db: db,
	}
}

func New() *DB {
	once.Do(initDB)

	return data
}
