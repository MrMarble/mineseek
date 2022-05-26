package main

import (
	"fmt"
	"net"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
	"github.com/mrmarble/mineseek/libs/database"
	"github.com/mrmarble/mineseek/libs/minecraft"
	//"github.com/mrmarble/mineseek/queue"
)

func main() {
	//queue := queue.New("ping", "servers")
	//
	/* 	queue.StartConsuming(func(s string) error {

		err = db.InsertSLP(mc)
		if err != nil {
			log.Printf("Error inserting SLP %v", err)
		}
		return nil
	}) */

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

func ping(addr string) (*minecraft.ServerListPing, error) {
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
	db := database.New()
	err = db.InsertSLP(mc)

	if err != nil {
		return nil, fmt.Errorf("Error inserting SLP %v", err)
	}
	return mc, nil
}
