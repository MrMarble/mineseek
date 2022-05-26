package minecraft

import (
	"encoding/json"

	"github.com/xrjr/mcutils/pkg/query"
)

type FullQuery map[string]interface{}

func newFullQuery(addr string, port int, query query.FullStat) (*FullQuery, error) {
	var slp FullQuery
	data, err := json.Marshal(query)
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(data, &slp)
	if err != nil {
		return nil, err
	}
	slp["address"] = addr
	slp["port"] = port
	return &slp, nil
}
func Query(addr string, port int) (*FullQuery, error) {
	stat, err := query.QueryFull(addr, port)
	if err != nil {
		return nil, err
	}
	return newFullQuery(addr, port, stat)
}
