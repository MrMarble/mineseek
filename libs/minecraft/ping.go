package minecraft

import "github.com/xrjr/mcutils/pkg/ping"

type ServerListPing struct {
	Version    string `json:"version"`
	MOTD       string `json:"motd"`
	Host       string `json:"host"`
	Favicon    string `json:"favicon"`
	Port       int    `json:"port"`
	MaxPlayers int    `json:"maxPlayers"`
}

func newSLPFromPing(host string, port int, ping ping.Infos) *ServerListPing {
	return &ServerListPing{
		Version:    ping.Version.Name,
		MOTD:       ping.Description,
		Favicon:    ping.Favicon,
		Host:       host,
		Port:       port,
		MaxPlayers: ping.Players.Max,
	}
}

func newSLPFromLegacyPing(host string, port int, ping ping.LegacyPingInfos) *ServerListPing {
	return &ServerListPing{
		Version:    ping.MinecraftVersion,
		MOTD:       ping.MOTD,
		Favicon:    "",
		Host:       host,
		Port:       port,
		MaxPlayers: ping.MaxPlayers,
	}
}

func Ping(addr string, port int) (*ServerListPing, error) {
	properties, _, err := ping.Ping(addr, port)

	if err == ping.ErrMalformedPacket {
		return pingLegacy(addr, port)
	}

	return newSLPFromPing(addr, port, properties.Infos()), err
}

func pingLegacy(addr string, port int) (*ServerListPing, error) {
	properties, _, err := ping.PingLegacy1_6_4(addr, port)

	if err == ping.ErrMalformedPacket {
		return pingOld(addr, port)
	}

	return newSLPFromLegacyPing(addr, port, properties), err
}

func pingOld(addr string, port int) (*ServerListPing, error) {
	properties, _, err := ping.PingLegacy(addr, port)

	if err == ping.ErrMalformedPacket {
		return pingLegacy(addr, port)
	}

	return newSLPFromLegacyPing(addr, port, properties), err
}
