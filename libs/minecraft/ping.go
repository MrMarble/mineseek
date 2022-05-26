package minecraft

import (
	"encoding/json"

	"github.com/xrjr/mcutils/pkg/ping"
)

type SLP map[string]interface{}
type input interface {
	ping.JSON | ping.LegacyPingInfos
}

func toSLP[T input](addr string, port int, p T) (*SLP, error) {
	var slp SLP
	data, err := json.Marshal(p)
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

// Ping automatically detects server version
func Ping(addr string, port int) (*SLP, error) {
	properties, _, err := ping.Ping(addr, port)
	if err == ping.ErrInvalidPacketType {
		return pingLegacy(addr, port)
	}

	if err != nil {
		return nil, err
	}
	return toSLP(addr, port, properties)
}

func pingLegacy(addr string, port int) (*SLP, error) {
	properties, _, err := ping.PingLegacy(addr, port)

	if err == ping.ErrInvalidPacketType {
		return pingOld(addr, port)
	}

	if err != nil {
		return nil, err
	}
	return toSLP(addr, port, properties)
}

func pingOld(addr string, port int) (*SLP, error) {
	properties, _, err := ping.PingLegacy1_6_4(addr, port)
	if err != nil {
		return nil, err
	}

	return toSLP(addr, port, properties)
}
