package ant

import (
    "strconv"
    "strings"
    "errors"
)

// https://zh.wikipedia.org/wiki/IPv4
func IP2Decimal(ip string)(n int64, err error){
    var ss []string = strings.Split(ip, ".")
    var b string
    var s string
    var i int64
    if len(ss) != 4{
        err = errors.New("IP Invalid")
        return
    }
    for _, s = range ss{
        i, err = strconv.ParseInt(s, 10, 64)
        if err != nil{
            return
        }
        s = strconv.FormatInt(i, 2)
        var j int
        var need int = 8 - len(s)
        for j=0;j<need;j++{
            s = "0" + s
        }
        b += s
    }
    n, _ = strconv.ParseInt(b, 2, 64)
    return
}

// https://zh.wikipedia.org/wiki/IPv4
func Decimal2IP(n int64)(ip string, err error){
    var ips []string = make([]string, 4)
    var b string
    var i int64
    b = strconv.FormatInt(n, 2)
    var need int = 32 - len(b)
    var j int
    for j=0;j<need;j++{
        b = "0" + b
    }
    i, _ = strconv.ParseInt(b[0:8], 2, 64)
    ips[0] = strconv.FormatInt(i, 10)
    i, _ = strconv.ParseInt(b[8:16], 2, 64)
    ips[1] = strconv.FormatInt(i, 10)
    i, _ = strconv.ParseInt(b[16:24], 2, 64)
    ips[2] = strconv.FormatInt(i, 10)
    i, _ = strconv.ParseInt(b[24:32], 2, 64)
    ips[3] = strconv.FormatInt(i, 10)
    ip = strings.Join(ips, ".")
    return
}

type CIDRInfo struct{
    First string
    Last string
    Block int64
    Network string
    Count int64
}

// wiki: http://goo.gl/AEUIi8
func CIDR(cidr string)(c *CIDRInfo, err error){
    c = new(CIDRInfo)
    var cs []string = strings.Split(cidr, "/")
    if len(cs) != 2{
        err = errors.New("CIDR Invalid")
        return
    }
    var ipd int64
    ipd, err = IP2Decimal(cs[0])
    if err != nil {
        return
    }
    var ipb string
    ipb = strconv.FormatInt(ipd, 2)
    var need int = 32 - len(ipb)
    var j int
    for j=0;j<need;j++{
        ipb = "0" + ipb
    }

    var n int64
    n, err = strconv.ParseInt(cs[1], 10, 64)
    if err != nil {
        return
    }
    if n<0 || n>32{
        err = errors.New("CIDR Invalid")
        return
    }
    c.Block = n

    var network string
    var networkI int64
    for j=0;j<int(n);j++{
        network += "1"
    }
    for j=0;j<32-int(n);j++{
        network += "0"
    }
    networkI, _ = strconv.ParseInt(network, 2, 64)
    network, _ = Decimal2IP(networkI)
    c.Network = network

    var first string = ipb[0:n]
    var firstI int64
    for j=0;j<32-int(n);j++{
        first = first + "0"
    }
    firstI, _ = strconv.ParseInt(first, 2, 64)
    first, _ = Decimal2IP(firstI)
    c.First = first

    var last string = ipb[0:n]
    var lastI int64
    for j=0;j<32-int(n);j++{
        last = last + "1"
    }
    lastI, _ = strconv.ParseInt(last, 2, 64)
    last, _ = Decimal2IP(lastI)
    c.Last = last

    c.Count = lastI - firstI + 1
    return
}
