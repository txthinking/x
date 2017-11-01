package ant

import (
	"testing"
)

func TestSHA1(t *testing.T) {
	s := "tx"
	s = SHA1(s)
	t.Log(s)
}

func TestSHA255(t *testing.T) {
	s := "tx"
	s = SHA256String(s)
	t.Log(s)
}
