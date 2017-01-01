package ant

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"errors"
	"io"
)

// AES256KeyLength is the length of key for AES 256 crypt
const AES256KeyLength = 32

// AESMake256Key cut or append empty data on the key
// and make sure the key lenth equal 32
func AESMake256Key(k []byte) []byte {
	if len(k) < AES256KeyLength {
		a := make([]byte, AES256KeyLength-len(k))
		return append(k, a...)
	}
	if len(k) > AES256KeyLength {
		return k[:AES256KeyLength]
	}
	return k
}

// AESEncrypt encrypt s with given k
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

// AESDecrypt decrypt c with given k
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
