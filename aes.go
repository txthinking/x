package ant

import(
    "crypto/aes"
    "crypto/cipher"
    "crypto/rand"
    "io"
    "encoding/hex"
    "errors"
    "bytes"
)

const AES_256_KEY_LENTH = 32

// Cut or append empty data on the key.
// make the key length equal 32
func AESMake256Key(key string)(k []byte){
    k = bytes.NewBufferString(key).Bytes()
    if len(k) < AES_256_KEY_LENTH{
        var a []byte = make([]byte, AES_256_KEY_LENTH-len(k))
        k = append(k, a...)
        return
    }
    if len(k) > AES_256_KEY_LENTH{
        k = k[:AES_256_KEY_LENTH]
        return
    }
    return
}

func AESEncrypt(s string, key string)(c string, err error){
    var k []byte
    var block cipher.Block
    var cfb cipher.Stream
    var cb []byte
    var iv []byte

    k = AESMake256Key(key)

    block, err = aes.NewCipher(k)
    if err != nil {
        return
    }

    cb = make([]byte, aes.BlockSize + len(s))
    iv = cb[:aes.BlockSize]
    _, err = io.ReadFull(rand.Reader, iv)
    if err != nil {
        return
    }

    cfb = cipher.NewCFBEncrypter(block, iv)
    cfb.XORKeyStream(cb[aes.BlockSize:], bytes.NewBufferString(s).Bytes())

    c = hex.EncodeToString(cb)
    return
}

func AESDecrypt(c string, key string)(s string, err error){
    var k []byte
    var block cipher.Block
    var cfb cipher.Stream
    var cb []byte
    var iv []byte

    k = AESMake256Key(key)

    cb, err = hex.DecodeString(c)
    if err != nil {
        return
    }

    block, err = aes.NewCipher(k)
    if err != nil {
        return
    }
    if len(cb) < aes.BlockSize{
        err = errors.New("crypt string is too short")
        return
    }

    iv = cb[:aes.BlockSize]
    cb = cb[aes.BlockSize:]

    cfb = cipher.NewCFBDecrypter(block, iv)
    cfb.XORKeyStream(cb, cb)
    s = bytes.NewBuffer(cb).String()
    return
}
