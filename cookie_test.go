package x

import "testing"

func TestCryptKV(t *testing.T) {
	kv := &CryptKV{
		LifeCycle: 7 * 24 * 60 * 60,
		AESKey:    []byte("qwertyuiopasdfghjklzqwertyuiopqw"),
	}
	k := "t"

	v := "hello"
	c, err := kv.Encrypt(k, v)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(c)
	v1, err := kv.Decrypt(c, k)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(v1)
}
