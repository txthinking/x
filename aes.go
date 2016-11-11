package ant

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"errors"
	"io"
)

const AES_256_KEY_LENTH = 32

// Cut or append empty data on the key.
// make the key length equal 32
func AESMake256Key(k []byte) []byte {
	if len(k) < AES_256_KEY_LENTH {
		var a []byte = make([]byte, AES_256_KEY_LENTH-len(k))
		return append(k, a...)
	}
	if len(k) > AES_256_KEY_LENTH {
		return k[:AES_256_KEY_LENTH]
	}
	return k
}

// aes 256 cfb
func AESEncrypt(s, k []byte) ([]byte, error) {
	block, err := aes.NewCipher(AESMake256Key(k))
	if err != nil {
		return nil, err
	}

	cb := make([]byte, aes.BlockSize+len(s))
	iv := cb[:aes.BlockSize]
	_, err = io.ReadFull(rand.Reader, iv)
	if err != nil {
		return nil, err
	}

	cfb := cipher.NewCFBEncrypter(block, iv)
	cfb.XORKeyStream(cb[aes.BlockSize:], s)
	return cb, nil
}

// aes 256 cfb
func AESDecrypt(c, k []byte) ([]byte, error) {
	block, err := aes.NewCipher(AESMake256Key(k))
	if err != nil {
		return nil, err
	}
	if len(c) < aes.BlockSize {
		err := errors.New("crypt data is too short")
		return nil, err
	}

	iv := c[:aes.BlockSize]
	cb := c[aes.BlockSize:]

	cfb := cipher.NewCFBDecrypter(block, iv)
	cfb.XORKeyStream(cb, cb)
	return cb, nil
}
