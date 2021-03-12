package x

import (
	"net"
)

// Resolver is a common interface for resolving URL's with alternate schemes
type Resolver interface {
	ResolveTCPAddr(network, addr string) (net.Addr, error)
	//(*net.TCPAddr, error)
	ResolveUDPAddr(network, addr string) (net.Addr, error)
	//(*net.UDPAddr, error)
}

type Resolve struct {
}

// DefaultDial is the default dialer in net package
var DefaultResolve = &Resolve{}

func (r *Resolve) ResolveTCPAddrS(network, addr string) (*net.TCPAddr, error) {
	return net.ResolveTCPAddr(network, addr)
}

func (r *Resolve) ResolveUDPAddrS(network, addr string) (*net.UDPAddr, error) {
	return net.ResolveUDPAddr(network, addr)
}

func (r *Resolve) ResolveTCPAddr(network, addr string) (net.Addr, error) {
	return r.ResolveTCPAddrS(network, addr)
}

func (r *Resolve) ResolveUDPAddr(network, addr string) (net.Addr, error) {
	return r.ResolveUDPAddrS(network, addr)
}
