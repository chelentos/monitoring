package validation

import (
	"fmt"
	"net"
	"strconv"
)

func ParseInt(value string, name string) (int, error) {
	num, err := strconv.Atoi(value)

	if err != nil {
		return 0, fmt.Errorf("%s must be int, got %s", name, value)
	}

	return num, nil
}

func CheckIntRange(value int, min int, max int, name string) error {
	if value < min || value > max {
		return fmt.Errorf("%s %d must be in %d – %d range", name, value, min, max)
	}

	return nil
}

func CheckAddr(addr string) error {
	_, port, err := net.SplitHostPort(addr)

	if err != nil {
		return fmt.Errorf("Invalid addr %q: %w", addr, err)
	}

	portNum, err := ParseInt(port, "port")

	if err != nil {
		return err
	}

	err = CheckIntRange(portNum, 1, 65535, "port")

	if err != nil {
		return err
	}

	return nil
}
