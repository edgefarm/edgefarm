package validate

import (
	"fmt"
	"net"
	"time"
)

func CheckFreePort(port int) bool {
	timeout := time.Second
	conn, err := net.DialTimeout("tcp", net.JoinHostPort("localhost", fmt.Sprintf("%d", port)), timeout)
	if err != nil {
		return true
	}
	if conn != nil {
		conn.Close()
		return false
	}
	return true
}
