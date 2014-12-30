package hecktest_test

import (
	"net"
	"testing"

	"github.com/theckman/hecktest"
	. "gopkg.in/check.v1"
)

func Test(t *testing.T) { TestingT(t) }

type TestSuite struct{}

var _ = Suite(&TestSuite{})

func (s *TestSuite) TestNewUDPListener(c *C) {
	var conn *net.UDPConn
	var err error

	conn, err = hecktest.NewUDPListener(8080)
	c.Assert(err, IsNil)

	defer conn.Close()

	c.Assert(conn.LocalAddr().String(), Equals, "127.0.0.1:8080")
}

func (s *TestSuite) TestUDPListen(c *C) {
	conn, err := hecktest.NewUDPListener(10101)
	c.Assert(err, IsNil)

	ctrl := make(chan int)
	outc := make(chan []byte)

	// shut down Goroutine below
	// the order is important here...
	defer close(ctrl)
	defer conn.Close()

	go hecktest.UDPListen(conn, ctrl, outc)

	udp, err := net.Dial("udp", "127.0.0.1:10101")

	c.Assert(err, IsNil)

	b := []byte("replace with witty joke")

	udp.Write(b)

	val, ok := <-outc
	c.Assert(ok, Equals, true)
	c.Assert(string(val), Equals, string(b))
}
