package x

import (
	"net"
)

// Dialer is a common interface for dialing
type Dialer interface {
	Dial(network, addr string) (net.Conn, error)
	DialTCP(network string, laddr, raddr net.Addr) (net.Conn, error)
	DialUDP(network string, laddr, raddr net.Addr) (net.PacketConn, error)
}

type Dial struct {
}

// DefaultDial is the default dialer in net package
var DefaultDial = &Dial{}

func (d *Dial) Dial(network, addr string) (net.Conn, error) {
	return net.Dial(network, addr)
}

<<<<<<< HEAD
func (d *Dial) DialTCPS(network string, laddr, raddr net.Addr) (*net.TCPConn, error) {
	return net.DialTCP(network, laddr.(*net.TCPAddr), raddr.(*net.TCPAddr))
}

func (d *Dial) DialUDPS(network string, laddr, raddr net.Addr) (*net.UDPConn, error) {
	return net.DialUDP(network, laddr.(*net.UDPAddr), raddr.(*net.UDPAddr))
}

func (d *Dial) DialTCP(network string, laddr, raddr net.Addr) (net.Conn, error) {
	return d.DialTCPS(network, laddr, raddr)
}

func (d *Dial) DialUDP(network string, laddr, raddr net.Addr) (net.PacketConn, error) {
	return d.DialUDPS(network, laddr, raddr)
}
