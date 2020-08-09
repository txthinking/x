package x

import "sync"

func NewBytesPool(n int) sync.Pool {
	return sync.Pool{
		New: func() interface{} {
			return make([]byte, n)
		},
	}
}

var BP65535 = NewBytesPool(65535)
var BP2064 = NewBytesPool(2064)
var BP2048 = NewBytesPool(2048)
var BP40 = NewBytesPool(40)
var BP20 = NewBytesPool(20)
var BP16 = NewBytesPool(16)
var BP12 = NewBytesPool(12)
var BP4 = NewBytesPool(4)
var BP2 = NewBytesPool(2)
