package ant

import(
    "crypto/md5"
    "encoding/hex"
    "io"
    "hash"
)

func MD5(s string)(r string){
    var h hash.Hash
    h = md5.New()
    io.WriteString(h, s)
    r = hex.EncodeToString(h.Sum(nil))
    return
}
