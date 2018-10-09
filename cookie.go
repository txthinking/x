package ant

import (
	"encoding/hex"
	"encoding/json"
	"errors"
	"net/http"
	"time"
)

// CryptKV can be used to crypt value by key and give it a lifecycle (s).
// Can be used in cookies
type CryptKV struct {
	AESKey    []byte
	LifeCycle int
}

// Encrypt key, value
func (kv *CryptKV) Encrypt(k string, v string) (string, error) {
	m := map[string]interface{}{
		"k":          k,
		"v":          v,
		"expired_at": time.Now().Add(time.Second * time.Duration(kv.LifeCycle)).Unix(),
	}
	b, err := json.Marshal(m)
	if err != nil {
		return "", err
	}
	b, err = AESCFBEncrypt(b, AESMake256Key(kv.AESKey))
	if err != nil {
		return "", err
	}
	return hex.EncodeToString(b), nil
}

// Decrypt key, value
func (kv *CryptKV) Decrypt(c string, k string) (string, error) {
	b, err := hex.DecodeString(c)
	if err != nil {
		return "", err
	}
	m := make(map[string]interface{})
	d, err := AESCFBDecrypt(b, AESMake256Key(kv.AESKey))
	if err != nil {
		return "", err
	}
	if err := json.Unmarshal(d, &m); err != nil {
		return "", err
	}
	if time.Unix(int64(m["expired_at"].(float64)), 0).Before(time.Now()) {
		return "", errors.New("Expired")
	}
	if m["k"].(string) != k {
		return "", errors.New("Unmarch key")
	}
	return m["v"].(string), nil
}

// Decrypt from http.Request
func (kv *CryptKV) DecryptFromRequest(r *http.Request, k string) (string, error) {
	c, err := r.Cookie(k)
	if err != nil {
		return "", err
	}
	v, err := kv.Decrypt(c.Value, k)
	if err != nil {
		return "", err
	}
	return v, nil
}
