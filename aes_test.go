package ant

import(
    "testing"
)

const KEY = "fuck, areyoufuckingwithme"

func TestAESMake256Key(t *testing.T){
    k := AESMake256Key(KEY)
    if len(k) != 32{
        t.Fatal("Not 32 length")
    }
}

func TestAESEncrypt(t *testing.T){
    s := "You know Nothing"
    r, err := AESEncrypt(s, KEY)
    if err != nil{
        t.Fatal(err)
    }
    t.Log(s, "<-->", r)
}

func TestAESDecrypt(t *testing.T){
    c := "94cedef55498b50b466cf70f460672f6ea83ace953ddb84999c31eb55f6301aa"
    r, err := AESDecrypt(c, KEY)
    if err != nil{
        t.Fatal(err)
    }
    t.Log(c, "<--->", r)
}
