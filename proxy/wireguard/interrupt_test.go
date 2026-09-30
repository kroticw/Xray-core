package wireguard

import (
	"io"
	"testing"
	"time"

	"github.com/xtls/xray-core/common"
	"github.com/xtls/xray-core/common/net"
)

// Outbounds end idle sessions by interrupting the link reader; for the
// WireGuard inbound that must close the flow so a blocked read returns.
func TestConnReaderInterruptClosesUDPFlow(t *testing.T) {
	m := &udpManager{handler: func(net.Conn, net.Destination) {}, m: make(map[string]*udpConn)}
	src := net.UDPDestination(net.ParseAddress("10.66.0.2"), 40000)
	dst := net.UDPDestination(net.ParseAddress("1.1.1.1"), 53)
	m.feed(src, dst, []byte{1})
	uc := m.m[src.NetAddr()]
	uc.ReadMultiBuffer() // drain the first packet

	r := newConnReader(uc)
	done := make(chan error, 1)
	go func() {
		_, err := r.ReadMultiBuffer()
		done <- err
	}()
	time.Sleep(50 * time.Millisecond)
	common.Interrupt(r)
	select {
	case err := <-done:
		if err != io.EOF {
			t.Fatalf("blocked read returned %v, want EOF", err)
		}
	case <-time.After(time.Second):
		t.Fatal("interrupt did not unblock the reader")
	}
}
