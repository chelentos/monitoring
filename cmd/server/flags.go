package main

import (
	"flag"
	"fmt"
	"net"
	"strconv"
)

var flagRunAddr string

func parseServerFlags() error {
	flag.StringVar(&flagRunAddr, "a", "localhost:8080", "address and port to run server")

	flag.Parse()

	err := checkAddr(flagRunAddr)

	return err
}

func checkAddr(addr string) error {
	_, port, err := net.SplitHostPort(addr)

	if err != nil {
		return fmt.Errorf("Invalid addr %q: %w", addr, err)
	}

	portNum, err := strconv.Atoi(port)

	if err != nil {
		return fmt.Errorf("Port must be int")
	}

	if portNum < 1 || portNum > 65535 {
		return fmt.Errorf("Port must be in 1 – 65535 range")
	}

	return nil
}
