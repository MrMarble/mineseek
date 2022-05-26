package main

import (
	"flag"
	"fmt"
	"net"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/mrmarble/mineseek/libs/database"
	"github.com/mrmarble/mineseek/libs/minecraft"
	"github.com/mrmarble/mineseek/libs/queue"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

var (
	useDB *bool
)

func init() {
	useDB = flag.Bool("dry", false, "save to database")
	flag.Parse()
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr, TimeFormat: time.RFC3339})
	log.Info().Bool("dry-run", *useDB).Msg("Service started")
}

func main() {
	queue := queue.New("ping", "servers")

	queue.StartConsuming(func(s string) error {
		_, err := ping(s)
		return err
	})

	e := echo.New()
	e.GET("/ping", func(c echo.Context) error {
		mc, err := ping(c.QueryParam("address"))
		if err != nil {
			return c.String(http.StatusNotFound, err.Error())
		}

		return c.JSON(http.StatusOK, mc)
	})
	e.Logger.Fatal(e.Start(":8080"))
}

func ping(addr string) (minecraft.SLP, error) {
	log.Info().Str("addr", addr).Msg("New address")
	host, port, err := net.SplitHostPort(addr)
	if err != nil {
		if _, ok := err.(*net.AddrError); ok {
			return ping(addr + ":25565")
		}
		return nil, fmt.Errorf("Error parsing address %v", err)
	}

	pint, err := strconv.Atoi(port)
	if err != nil {
		return nil, fmt.Errorf("Error parsing port %v", err)
	}

	mc, err := minecraft.Ping(host, pint)
	if err != nil {
		return nil, fmt.Errorf("Error pinging server %v", err)
	}

	if !*useDB {
		db := database.New()
		defer db.Disconnect()
		err = db.InsertSLP(mc)

		if err != nil {
			return nil, fmt.Errorf("Error inserting SLP %v", err)
		}

	}
	return mc, nil
}
