package ant

import (
	"bytes"
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

// AESCFBEncrypt encrypt s with given k
func AESCFBEncrypt(s, k []byte) ([]byte, error) {
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
func AESCFBDecrypt(c, k []byte) ([]byte, error) {
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

// AESCBCEncrypt encrypt s with given k
func AESCBCEncrypt(s, k []byte) ([]byte, error) {
	if len(s)%aes.BlockSize != 0 {
		return nil, errors.New("invalid length of s")
	}
	block, err := aes.NewCipher(AESMake256Key(k))
	if err != nil {
		return nil, err
	}
	cb := make([]byte, aes.BlockSize+len(s))
	iv := cb[:aes.BlockSize]
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		return nil, err
	}
	mode := cipher.NewCBCEncrypter(block, iv)
	mode.CryptBlocks(cb[aes.BlockSize:], s)
	return cb, nil
}

// AESCBCDecrypt decrypt c with given k
func AESCBCDecrypt(c, k []byte) ([]byte, error) {
	if len(c) < aes.BlockSize {
		return nil, errors.New("c too short")
	}
	block, err := aes.NewCipher(AESMake256Key(k))
	if err != nil {
		return nil, err
	}

	iv := c[:aes.BlockSize]
	cb := c[aes.BlockSize:]

	mode := cipher.NewCBCDecrypter(block, iv)
	mode.CryptBlocks(cb, cb)
	return cb, nil
}

// PKCS5 padding
func PKCS5Padding(c []byte, blockSize int) []byte {
	pl := blockSize - len(c)%blockSize
	p := bytes.Repeat([]byte{byte(pl)}, pl)
	return append(c, p...)
}

// PKCS5 unpadding
func PKCS5UnPadding(s []byte) ([]byte, error) {
	l := len(s)
	if l == 0 {
		return nil, errors.New("s too short")
	}
	pl := int(s[l-1])
	if l < pl {
		return nil, errors.New("s too short")
	}
	return s[:(l - pl)], nil
}
