package scanner

import (
	"encoding/binary"
	"fmt"
	"net"

	"github.com/rs/zerolog/log"
)

// createHostRange converts a input ip addr string to a slice of ips on the cidr.
func createHostRange(netw string) []string {
	_, ipv4Net, err := net.ParseCIDR(netw)
	if err != nil {
		log.Fatal().Err(err).Msg("failed to parse cidr")
	}

	mask := binary.BigEndian.Uint32(ipv4Net.Mask)
	start := binary.BigEndian.Uint32(ipv4Net.IP)
	finish := (start & mask) | (mask ^ 0xffffffff)

	var hosts []string

	for i := start; i <= finish; i++ {
		ip := make(net.IP, 4)
		binary.BigEndian.PutUint32(ip, i)
		hosts = append(hosts, ip.String())
	}

	return hosts
}

// getLocalRange returns local ip range or defaults on error to most common.
func getLocalIP() (string, error) {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		return "", err
	}

	for _, address := range addrs {
		// check the address type and if it is not a loopback the display it
		if ipnet, ok := address.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				return ipnet.IP.String(), err
			}
		}
	}

	return "", fmt.Errorf("No IP Found")
}

func canSocketBind(laddr string) bool {
	// Check if user can listen on socket
	listenAddr, err := net.ResolveIPAddr("ip4", laddr)
	if err != nil {
		return false
	}

	conn, err := net.ListenIP("ip4:tcp", listenAddr)
	if err != nil {
		return false
	}

	conn.Close()

	return true
}
