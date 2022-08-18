package scanner

import (
	"fmt"
	"net"
	"strconv"
	"sync"
	"time"

	"github.com/rs/zerolog/log"
	"github.com/xrjr/mcutils/pkg/ping"

	"github.com/mrmarble/mineseek/internal/minecraft"
)

type Options struct {
	Host     string
	Protocol string
	Ports    []int
	Stealth  bool
	Rate     int
	Timeout  time.Duration
	SLP      bool
}

type Result struct {
	Host    string
	Port    int
	Open    bool
	SLP     *ping.JSON
	Latency int
}

type ResultChan chan Result

func Scan(opts *Options) (ResultChan, error) {
	laddr, err := getLocalIP()
	if err != nil {
		return nil, err
	}

	if opts.Stealth {
		if !canSocketBind(laddr) {
			return nil, fmt.Errorf("socket: operation not permitted")
		}
	}

	return scan(laddr, opts)
}

func scan(laddr string, opts *Options) (ResultChan, error) {
	cidr := createHostRange(opts.Host)
	log.Debug().Str("cidr", opts.Host).Int("Hosts", len(cidr)).Msg("scanning host range")
	tasks := len(cidr) * len(opts.Ports)
	hostsChan := make(chan string)
	resultsChan := make(chan Result, tasks)

	var wg sync.WaitGroup

	wg.Add(opts.Rate)

	go func() {
		for _, host := range cidr {
			hostsChan <- host
		}

		close(hostsChan)
	}()

	go func() {
		wg.Wait()
		close(resultsChan)
	}()

	worker := func(id int) {
		log.Debug().Int("id", id).Msg("Worker started")

		for {
			host, ok := <-hostsChan
			if !ok {
				log.Debug().Int("id", id).Msg("Worker exiting")
				wg.Done()

				return
			}

			if host != "" {
				for _, port := range opts.Ports {
					log.Debug().Int("Worker", id).Str("host", host).Int("port", port).Bool("Stealth", opts.Stealth).Msg("Scanning")

					if opts.Stealth {
						scanPortSyn(resultsChan, opts.Protocol, host, port, laddr, opts.Timeout, opts.SLP)
					} else {
						scanPort(resultsChan, port, host, opts.Protocol, opts.Timeout, opts.SLP)
					}
				}
			}
		}
	}

	for i := 0; i < opts.Rate; i++ {
		log.Debug().Int("id", i+1).Msg("Starting worker")

		go worker(i + 1)
	}

	return resultsChan, nil
}

// scanPort scans a single ip port combo
// This detection method only works on some types of services
// but is a reasonable solution for this application.
func scanPort(resultChannel chan<- Result, port int, host, protocol string, timeout time.Duration, slp bool) {
	result := Result{
		Host: host,
		Port: port,
	}
	address := host + ":" + strconv.Itoa(port)

	conn, err := net.DialTimeout(protocol, address, timeout)
	if err != nil {
		result.Open = false
		resultChannel <- result

		return
	}

	defer conn.Close()

	result.Open = true

	if slp {
		properties, latency, err := minecraft.Ping(host, port, timeout)
		if err != nil {
			log.Debug().Err(err).Msg("Error pinging server")
		}

		result.SLP = &properties
		result.Latency = latency
	}

	resultChannel <- result
}

func scanPortSyn(resultChannel chan<- Result, protocol, host string, port int,
	laddr string, timeout time.Duration, slp bool,
) {
	result := Result{
		Host: host,
		Port: port,
	}
	ack := make(chan bool, 1)

	go recvSynAck(laddr, host, uint16(port), ack)
	sendSyn(laddr, host, uint16(random(10000, 65535)), uint16(port))

	select {
	case r := <-ack:
		result.Open = r

		if r && slp {
			properties, latency, err := minecraft.Ping(host, port, timeout)
			if err == nil {
				result.SLP = &properties
				result.Latency = latency
			}
		}
		resultChannel <- result

		return
	case <-time.After(timeout):
		result.Open = false
		resultChannel <- result

		return
	}
}
