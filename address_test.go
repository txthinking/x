package x

import (
	"testing"
)

func TestHasPort(t *testing.T) {
	if HasPort("1.1.1.1") {
		t.Fatal("No Port")
	}
	if !HasPort("1.1.1.1:80") {
		t.Fatal("Has Port")
	}
}
