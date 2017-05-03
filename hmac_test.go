package ant

import (
	"encoding/hex"
	"testing"
)

func TestHmacSha256(t *testing.T) {
	s := []byte("txthinking")
	r, err := HmacSha256(s, []byte(KEY))
	if err != nil {
		t.Fatal(err)
	}
	t.Log(hex.EncodeToString(r), len(r))

	ok, err := CheckHmacSha256(s, r, []byte(KEY))
	if err != nil {
		t.Fatal(err)
	}
	if !ok {
		t.Fatal("HmacSha256 Error")
	}
}
