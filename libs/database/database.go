package database

import (
	"context"
	"log"
	"os"
	"sync"
	"time"

	"github.com/mrmarble/mineseek/libs/minecraft"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

var (
	data *DB
	once sync.Once
)

type DB struct {
	ctx    context.Context
	cancel context.CancelFunc
	client *mongo.Client
}

func (db *DB) Disconnect() {
	db.cancel()
	if err := db.client.Disconnect(db.ctx); err != nil {
		log.Fatal(err)
	}
}

func (db *DB) InsertSLP(slp *minecraft.SLP) error {
	coll := db.client.Database("mineseek").Collection("servers")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	doc, err := bson.Marshal(slp)
	if err != nil {
		return err
	}
	_, err = coll.InsertOne(ctx, doc)
	return err
}

func (db *DB) InsertQuery(query *minecraft.FullQuery) error {
	coll := db.client.Database("mineseek").Collection("queries")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	doc, err := bson.Marshal(query)
	if err != nil {
		return err
	}
	_, err = coll.InsertOne(ctx, doc)
	return err
}

func getConnection(ctx context.Context) (*mongo.Client, error) {
	return mongo.Connect(ctx, options.Client().ApplyURI(os.Getenv("DATABASE_URL")))
}

func innitDB() {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	db, err := getConnection(ctx)
	if err != nil {
		log.Panic(err)
	}

	data = &DB{
		ctx:    ctx,
		client: db,
		cancel: cancel,
	}
}

func New() *DB {
	if data != nil {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		err := data.client.Ping(ctx, readpref.Primary())
		if err != nil {
			innitDB()
			return data
		}
	} else {
		innitDB()
	}
	return data
}
