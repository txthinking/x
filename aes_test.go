package ant

import (
	"testing"
)

const KEY = "fuck, areyoufuckingwithme"

func TestAESMake256Key(t *testing.T) {
	k := AESMake256Key([]byte(KEY))
	if len(k) != 32 {
		t.Fatal("Not 32 length")
	}
}

func TestAES(t *testing.T) {
	s := []byte("You know Nothing")
	r, err := AESEncrypt(s, []byte(KEY))
	if err != nil {
		t.Fatal(err)
	}
	t.Log("input:", string(s))

	r, err = AESDecrypt(r, []byte(KEY))
	if err != nil {
		t.Fatal(err)
	}
	t.Log("output:", string(r))
	if string(s) != string(r) {
		t.Fatal("AES Error")
	}
}
