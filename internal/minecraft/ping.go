package minecraft

import (
	"errors"
	"time"

	"github.com/xrjr/mcutils/pkg/ping"
)

// Ping returns the server list ping infos (JSON-like object), and latency of a minecraft server.
func Ping(hostname string, port int, timeout time.Duration) (ping.JSON, int, error) {
	client := ping.NewClient(hostname, port)
	client.DialTimeout = timeout
	client.ReadTimeout = timeout

	err := client.Connect()
	if err != nil {
		return nil, -1, err
	}

	handshake, err := client.Handshake()
	if err != nil {
		return nil, -1, err
	}

	latency, err := client.Ping()

	// Some forge servers respond to ping request with the handshake response.
	// In this case, a ErrInvalidPacketType will be returned.
	// We'll be ingoring this error because it doesn't have any side effect, since :
	//   - we don't retrieve any information from the pong response packet
	//   - connection is closed right after
	if err != nil && !errors.Is(err, ping.ErrInvalidPacketType) {
		return nil, -1, err
	}

	err = client.Disconnect()
	if err != nil {
		return nil, -1, err
	}

	return handshake.Properties, latency, nil
}
