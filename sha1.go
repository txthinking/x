package ant

import(
    "crypto/sha1"
    "encoding/hex"
    "io"
    "hash"
)

func SHA1(s string)(r string){
    var h hash.Hash
    h = sha1.New()
    io.WriteString(h, s)
    r = hex.EncodeToString(h.Sum(nil))
    return
}
