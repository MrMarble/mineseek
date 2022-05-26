package database

import (
	"database/sql"
	"log"
	"os"
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
	_, err := db.db.Exec("INSERT INTO servers VALUES ($1, $2, $3, $4, $5, $6) ON CONFLICT DO UPDATE", slp.Host, slp.Port, slp.Version, slp.Favicon, slp.MOTD, slp.MaxPlayers)
	return err
}

func (db *DB) InsertQuery(query *minecraft.FullQuery) error {
	_, err := db.db.Exec("INSERT INTO queries VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9) ON CONFLICT DO UPDATE", query.Host, query.Port, query.GameType, query.GameID, query.Version, query.Plugins, query.Map, query.MaxPlayers)
	return err
}

func (db *DB) InsertPlayers(query *minecraft.FullQuery) error {
	tx, err := db.db.Begin()
	if err != nil {
		return err
	}
	for _, player := range query.Players {
		uuid, err := minecraft.GetUUID(player)
		if err != nil {
			return err
		}
		_, err = tx.Exec("INSERT INTO players VALUES ($1, $2) ON CONFLICT DO UPDATE", uuid, player)
		if err != nil {
			return err
		}
	}
	err = tx.Commit()
	return err
}

func getConnection() (*sql.DB, error) {
	return sql.Open("postgres", os.Getenv("DATABASE_URL"))
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
