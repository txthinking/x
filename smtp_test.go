package ant

import(
    "testing"
    "net/mail"
)

func TestMakeBoundary(t *testing.T){
    s := MakeBoundary()
    t.Log(s)
}

func TestChunkSplit(t *testing.T){
    s := `
aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa
bbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbb
ccccccccccccccccccccccccccccccccccccccc
`
    s, err := ChunkSplit(s)
    if err != nil{
        t.Fatal(err)
    }
    t.Log(s)
}

func TestSendMail(t *testing.T){
	p := ""
    s := &SMTP{
        "smtp.ym.163.com",
        25,
        "bot@ym.txthinking.com",
        p,
        false,
    }
    f := &mail.Address{
        Name: "雷锋",
        Address: "bot@ym.txthinking.com",
    }
    ts := []*mail.Address{
        &mail.Address{
            Name: "雷锋",
            Address: "cloud@txthinking.com",
        },
    }
    m := &Message{
        From: f,
        To: ts,
        Subject: "hello",
        Body: "哈哈",
        Att: []string{
            "/etc/hosts",
        },
    }

    err := s.Send(m)
    if err != nil{
        t.Fatal(err)
    }
}
