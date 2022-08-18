package utils

import (
	"fmt"
	"net"
	"strconv"
	"strings"
)

func ParsePorts(ports string) ([]int, error) {
	port, err := strconv.Atoi(ports)
	if err == nil {
		return []int{port}, nil
	}

	var portList []int

	if strings.ContainsRune(ports, ',') {
		for _, p := range strings.Split(ports, ",") {
			port, err := strconv.Atoi(p)
			if err != nil {
				return nil, err
			}

			portList = append(portList, port)
		}

		return portList, nil
	} else if strings.ContainsRune(ports, '-') {
		portRange := strings.Split(ports, "-")
		if len(portRange) != 2 {
			return nil, fmt.Errorf("invalid port range")
		}
		start, err := strconv.Atoi(portRange[0])
		if err != nil {
			return nil, err
		}
		end, err := strconv.Atoi(portRange[1])
		if err != nil {
			return nil, err
		}
		for i := start; i <= end; i++ {
			portList = append(portList, i)
		}

	}

	return nil, fmt.Errorf("invalid port range")
}

func AvailableHosts(cidr string) int {
	_, ipnet, err := net.ParseCIDR(cidr)
	if err != nil {
		return 0
	}

	ones, bits := ipnet.Mask.Size()

	return (1 << uint(bits-ones))
}
