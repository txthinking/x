package x

import (
	"net"
)

// Resolver is a common interface for dialing
type Resolver interface {
	Resolve(network, addr string) (net.Addr, error)
	ResolveTCPAddr(network, dest string) (net.Addr, error)
	ResolveUDPAddr(network, dest string) (net.Addr, error)
}

type Resolve struct {
}

// DefaultDial is the default dialer in net package
var DefaultResolve = &Resolve{}

func (d *Resolve) Resolve(network, addr string) (net.Addr, error) {
	if network == "tcp" || network == "tcp6" || network == "tcp4" {
		return net.ResolveTCPAddr(network, addr)
	} else if network == "udp" || network == "udp6" || network == "udp4" {
		return net.ResolveUDPAddr(network, addr)
	}
	return net.ResolveIPAddr(network, addr)
}

func (d *Resolve) ResolveTCPAddr(network, dest string) (net.Addr, error) {
	return net.ResolveTCPAddr("tcp", dest)
}

func (d *Resolve) ResolveUDPAddr(network, dest string) (net.Addr, error) {
	return net.ResolveUDPAddr("udp", dest)
}
