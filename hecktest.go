// Package hecktest is a shared library for
// code that can be reused for unit testing.
//
// hecktest is relesed under the 3-Clause BSD License.
//
// hecktest contains modifications of functions written for
// PagerDuty/godspeed, which is also under the 3-Clause BSD License.
package hecktest

import (
	"bytes"
	"fmt"
	"net"
)

// UDPListen is a helper function for running a test UDP server.
// It's meant to be ran as a Goroutine with each full UDP payload send
// through the out channel.
//
// To shut down this Goroutine you should invoke close the ctrl channel
// and then invoke conn.Close(). The ordering is important to avoid this
// Goroutine entering an infinite loop until you close the ctrl channel
func UDPListen(conn *net.UDPConn, ctrl <-chan int, out chan<- []byte) {
	for {
		select {
		case _, ok := <-ctrl:
			if !ok {
				close(out)
				return
			}
		default:
			buf := make([]byte, 8193)

			_, err := conn.Read(buf)

			// err could just be an early end
			// of stream or something
			if err != nil {
				continue
			}

			// trim NULL bytes and send it through
			out <- bytes.Trim(buf, "\x00")
		}
	}
}

// NewUDPListener is a function that returns a *net.UDPConn instance.
// It listens on localhost by default on the port provided as an argument.
func NewUDPListener(port uint16) (*net.UDPConn, error) {
	addr, err := net.ResolveUDPAddr("udp", fmt.Sprintf("127.0.0.1:%d", port))

	if err != nil {
		return nil, fmt.Errorf("Address resolution failure: %v", err.Error())
	}

	l, err := net.ListenUDP("udp", addr)

	if err != nil {
		return nil, fmt.Errorf("error: %v", err.Error())
	}

	return l, nil
}
