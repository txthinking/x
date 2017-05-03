package ant

import (
	"encoding/hex"
	"testing"
)

const KEY = "areyoufuckingwithmemustbe32bytes"

func TestAESMake256Key(t *testing.T) {
	k := AESMake256Key([]byte(KEY))
	if len(k) != 32 {
		t.Fatal("Not 32 length")
	}
}

func TestAESCFB(t *testing.T) {
	s := []byte("txthinking")
	r, err := AESCFBEncrypt(s, []byte(KEY))
	if err != nil {
		t.Fatal(err)
	}
	t.Log("input:", string(s))

	r, err = AESCFBDecrypt(r, []byte(KEY))
	if err != nil {
		t.Fatal(err)
	}
	t.Log(len(s), len(r))
	t.Log("output:", string(r))
	if string(s) != string(r) {
		t.Fatal("AES CFB Error")
	}
}

func TestAESCBC(t *testing.T) {
	s := []byte("txthinking")
	t.Log("input:", hex.EncodeToString(s), len(s))

	sp := PKCS5Padding(s, 16)

	r, err := AESCBCEncrypt(sp, []byte(KEY))
	if err != nil {
		t.Fatal(err)
	}
	t.Log("encrypt:", hex.EncodeToString(r), len(r))

	r, err = AESCBCDecrypt(r, []byte(KEY))
	if err != nil {
		t.Fatal(err)
	}
	r, err = PKCS5UnPadding(r)
	if err != nil {
		t.Fatal(err)
	}
	t.Log("output:", hex.EncodeToString(r), len(r))

	if string(s) != string(r) {
		t.Fatal("AES CBC Error")
	}
}
