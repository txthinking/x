package ant

import(
    "net/smtp"
    "time"
    "bytes"
    "io"
    "crypto/tls"
    "fmt"
    "encoding/base64"
    "strconv"
    "path/filepath"
    "math/rand"
    "io/ioutil"
    "net/mail"
)

type SMTP struct{
    Server string
    Port int
    UserName string
    Password string
    IsTLS bool
}

// Create boundary for MIME data.
func MakeBoundary()(b string){
    b = strconv.FormatInt(time.Now().UnixNano(), 10)
    b += strconv.FormatInt(rand.New(rand.NewSource(time.Now().UnixNano())).Int63(), 10)
    b = MD5(b)
    return
}

// Chunk data using RFC 2045.
func ChunkSplit(s string)(r string, err error){
    const LENTH = 76
    var bfr, bfw *bytes.Buffer
    var data, block []byte

    data = make([]byte, 0)
    block = make([]byte, LENTH)
    bfr = bytes.NewBufferString(s)
    bfw = bytes.NewBuffer(data)
    var l int
    for {
        l, err = bfr.Read(block)
        if err == io.EOF{
            err = nil
            break
        }
        if err != nil{
            return
        }
        _, err = bfw.Write(block[:l])
        if err != nil{
            return
        }
        _, err = bfw.WriteString("\r\n")
        if err != nil{
            return
        }
    }
    r = bfw.String()
    return
}

// RFC 821,822,1869,2821
// from/to can be "some@domain.com"
// or RFC 5322 address, e.g. "Barry Gibbs <bg@example.com>"
// NOTICE: If the name is not ASCII, you should be do some encoding,
// eg:"=?utf-8?B?"+base64.StdEncoding.EncodeToString([]byte("名字"))+"?="
// att is the file path for attachments or nil if no attachment.
func (m *SMTP)Send(msg *Message)(err error){
    var client *smtp.Client
    var auth smtp.Auth

    var i int

    client, err = smtp.Dial(m.Server + ":" + strconv.Itoa(m.Port))
    if err != nil {
        return
    }
    err = client.Hello(m.Server)
    if err != nil {
        return
    }
    if m.IsTLS{
        err = client.StartTLS(&tls.Config{ServerName: m.Server, InsecureSkipVerify: true})
        if err != nil {
            return
        }
    }
    auth = smtp.PlainAuth("", m.UserName, m.Password, m.Server)
    err = client.Auth(auth)
    if err != nil {
        return
    }
    err = client.Mail(msg.From.Address)
    if err != nil {
        return
    }
    for i, _ = range msg.To{
        err = client.Rcpt(msg.To[i].Address)
        if err != nil {
            return
        }
    }
    var in io.WriteCloser
    in, err = client.Data()
    if err != nil {
        return
    }
    /*
    var data string
    data, err = msg.String()
    if err != nil {
        return
    }
    _, err = io.WriteString(in, data)
    if err != nil {
        return
    }
    */
    var r io.Reader
    r, err = msg.Reader()
    if err != nil {
        return
    }
    _, err = io.Copy(in, r)
    if err != nil {
        return
    }
    err = in.Close()
    if err != nil {
        return
    }
    err = client.Quit()
    if err != nil {
        return
    }
    return
}

type Message struct{
    From *mail.Address
    To []*mail.Address
    Subject string
    Body string
    Att []string
}

// MIME RFC2045
// https://en.wikipedia.org/wiki/MIME
func (m *Message)String()(data string, err error){
    var s string
    var i int
    var bs []byte

    // prepare body data
    m.Body, err = ChunkSplit(base64.StdEncoding.EncodeToString(bytes.NewBufferString(m.Body).Bytes()))
    if err != nil {
        return
    }

    // prepare attachment data
    var attms []map[string]string
    var attm map[string]string
    if len(m.Att) != 0{
        attms = make([]map[string]string, len(m.Att))
        for i, s = range m.Att{
            attm = make(map[string]string)
            bs, err = ioutil.ReadFile(s)
            if err != nil{
                return
            }
            attm["name"] = filepath.Base(s)
            attm["data"], err = ChunkSplit(base64.StdEncoding.EncodeToString(bs))
            if err != nil {
                return
            }
            attms[i] = attm
        }
    }

    // prepare mail data
    var boundaryAlternative, boundaryMixed string

    data += fmt.Sprintf("Subject: %s\r\n", m.Subject)
    data += fmt.Sprintf("From: %s\r\n", m.From.String())
    data += fmt.Sprintf("MIME-Version: 1.0\r\n")
    s = ""
    for i, _ = range m.To{
        s += m.To[i].String() + ", "
    }
    data += fmt.Sprintf("To: %s\r\n", s[:len(s)-2])

    boundaryAlternative = MakeBoundary()
    if len(m.Att) != 0{
        boundaryMixed = MakeBoundary()
        data += fmt.Sprintf("Content-Type: multipart/mixed; boundary=\"%s\"\r\n\r\n", boundaryMixed)
        data += fmt.Sprintf("--%s\r\n", boundaryMixed)
    }

    data += fmt.Sprintf("Content-Type: multipart/alternative; boundary=\"%s\"\r\n\r\n", boundaryAlternative)
    data += fmt.Sprintf("--%s\r\n", boundaryAlternative)
    data += fmt.Sprintf("Content-Type: text/plain; charset=utf-8\r\n")
    data += fmt.Sprintf("Content-Transfer-Encoding: base64\r\n\r\n")
    data += fmt.Sprintf("%s\r\n\r\n", m.Body)
    data += fmt.Sprintf("--%s\r\n", boundaryAlternative)
    data += fmt.Sprintf("Content-Type: text/html; charset=utf-8\r\n")
    data += fmt.Sprintf("Content-Transfer-Encoding: base64\r\n\r\n")
    data += fmt.Sprintf("%s\r\n\r\n", m.Body)
    data += fmt.Sprintf("--%s--\r\n", boundaryAlternative)

    if len(m.Att) != 0{
        for _, attm = range attms{
            data += fmt.Sprintf("--%s\r\n", boundaryMixed)
            data += fmt.Sprintf("Content-Type: application/octet-stream; name=\"%s\"\r\n", attm["name"])
            data += fmt.Sprintf("Content-Transfer-Encoding: base64\r\n")
            data += fmt.Sprintf("Content-Disposition: attachment; filename=\"%s\"\r\n\r\n", attm["name"])
            data += fmt.Sprintf("%s\r\n\r\n", attm["data"])
        }
        data += fmt.Sprintf("--%s--\r\n", boundaryMixed)
    }
    return
}

// MIME RFC2045
// https://en.wikipedia.org/wiki/MIME
func (m *Message)Reader()(r io.Reader, err error){
    var s string
    var i int
    var bs []byte
    var bf *bytes.Buffer = bytes.NewBufferString("")

    // prepare body data
    m.Body, err = ChunkSplit(base64.StdEncoding.EncodeToString(bytes.NewBufferString(m.Body).Bytes()))
    if err != nil {
        return
    }

    // prepare attachment data
    var attms []map[string]string
    var attm map[string]string
    if len(m.Att) != 0{
        attms = make([]map[string]string, len(m.Att))
        for i, s = range m.Att{
            attm = make(map[string]string)
            bs, err = ioutil.ReadFile(s)
            if err != nil{
                return
            }
            attm["name"] = filepath.Base(s)
            attm["data"], err = ChunkSplit(base64.StdEncoding.EncodeToString(bs))
            if err != nil {
                return
            }
            attms[i] = attm
        }
    }

    // prepare mail data
    var boundaryAlternative, boundaryMixed string

    bf.WriteString(fmt.Sprintf("Subject: %s\r\n", m.Subject))
    bf.WriteString(fmt.Sprintf("From: %s\r\n", m.From.String()))
    bf.WriteString(fmt.Sprintf("MIME-Version: 1.0\r\n"))
    s = ""
    for i, _ = range m.To{
        s += m.To[i].String() + ", "
    }
    bf.WriteString(fmt.Sprintf("To: %s\r\n", s[:len(s)-2]))

    boundaryAlternative = MakeBoundary()
    if len(m.Att) != 0{
        boundaryMixed = MakeBoundary()
        bf.WriteString(fmt.Sprintf("Content-Type: multipart/mixed; boundary=\"%s\"\r\n\r\n", boundaryMixed))
        bf.WriteString(fmt.Sprintf("--%s\r\n", boundaryMixed))
    }

    bf.WriteString(fmt.Sprintf("Content-Type: multipart/alternative; boundary=\"%s\"\r\n\r\n", boundaryAlternative))
    bf.WriteString(fmt.Sprintf("--%s\r\n", boundaryAlternative))
    bf.WriteString(fmt.Sprintf("Content-Type: text/plain; charset=utf-8\r\n"))
    bf.WriteString(fmt.Sprintf("Content-Transfer-Encoding: base64\r\n\r\n"))
    bf.WriteString(fmt.Sprintf("%s\r\n\r\n", m.Body))
    bf.WriteString(fmt.Sprintf("--%s\r\n", boundaryAlternative))
    bf.WriteString(fmt.Sprintf("Content-Type: text/html; charset=utf-8\r\n"))
    bf.WriteString(fmt.Sprintf("Content-Transfer-Encoding: base64\r\n\r\n"))
    bf.WriteString(fmt.Sprintf("%s\r\n\r\n", m.Body))
    bf.WriteString(fmt.Sprintf("--%s--\r\n", boundaryAlternative))

    if len(m.Att) != 0{
        for _, attm = range attms{
            bf.WriteString(fmt.Sprintf("--%s\r\n", boundaryMixed))
            bf.WriteString(fmt.Sprintf("Content-Type: application/octet-stream; name=\"%s\"\r\n", attm["name"]))
            bf.WriteString(fmt.Sprintf("Content-Transfer-Encoding: base64\r\n"))
            bf.WriteString(fmt.Sprintf("Content-Disposition: attachment; filename=\"%s\"\r\n\r\n", attm["name"]))
            bf.WriteString(fmt.Sprintf("%s\r\n\r\n", attm["data"]))
        }
        bf.WriteString(fmt.Sprintf("--%s--\r\n", boundaryMixed))
    }
    r = bf
    return
}
