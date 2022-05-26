package minecraft

import "github.com/xrjr/mcutils/pkg/query"

type FullQuery struct {
	Hostname   string
	GameType   string
	GameID     string
	Version    string
	Map        string
	Host       string
	Port       int
	MaxPlayers string
	Players    []string
	Plugins    string
}

func newFullQuery(addr string, port int, query query.FullStat) *FullQuery {
	return &FullQuery{
		Hostname:   query.Properties["hostname"],
		GameType:   query.Properties["game type"],
		GameID:     query.Properties["game_id"],
		Version:    query.Properties["version"],
		Map:        query.Properties["map"],
		MaxPlayers: query.Properties["maxplayers"],
		Plugins:    query.Properties["plugins"],
		Players:    query.OnlinePlayers,
		Port:       port,
		Host:       addr,
	}
}

func Query(addr string, port int) (*FullQuery, error) {
	stat, err := query.QueryFull(addr, port)
	if err != nil {
		return nil, err
	}
	return newFullQuery(addr, port, stat), nil
}
