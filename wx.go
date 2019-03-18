package x

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"time"

	simplejson "github.com/bitly/go-simplejson"
)

func WxAppMakePrepayID(outTradeNo, totalFee int64, spbillCreateIP, body, notifyURL, appid, mchID, key string) (string, error) {
	nonceStr := strconv.FormatInt(RandomNumber(), 10)
	s := `
	<xml>
		<appid>%s</appid>
		<mch_id>%s</mch_id>
		<nonce_str>%s</nonce_str>
		<sign>%s</sign>
		<body>%s</body>
		<out_trade_no>%d</out_trade_no>
		<total_fee>%d</total_fee>
		<spbill_create_ip>%s</spbill_create_ip>
		<notify_url>%s</notify_url>
		<trade_type>APP</trade_type>
	</xml>
	`
	s0 := fmt.Sprintf("appid=%s&body=%s&mch_id=%s&nonce_str=%s&notify_url=%s&out_trade_no=%d&spbill_create_ip=%s&total_fee=%d&trade_type=APP&key=%s", appid, body, mchID, nonceStr, notifyURL, outTradeNo, spbillCreateIP, totalFee, key)
	sign := strings.ToUpper(MD5(s0))
	s = fmt.Sprintf(s, appid, mchID, nonceStr, sign, body, outTradeNo, totalFee, spbillCreateIP, notifyURL)

	c := &http.Client{
		Timeout: 10 * time.Second,
	}
	res, err := c.Post("https://api.mch.weixin.qq.com/pay/unifiedorder", "application/xml", bytes.NewReader([]byte(s)))
	if err != nil {
		return "", err
	}
	defer res.Body.Close()
	b, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return "", err
	}

	reg := regexp.MustCompile(`<(.+)>(.+)<`)
	ss := reg.FindAllStringSubmatch(strings.Replace(strings.Replace(string(b), "<xml>", "", -1), "</xml>", "", -1), -1)
	ss1 := make([]string, 0)
	var returnCode, returnMsg, resultCode, errCodeDes, prepayID string
	for _, v := range ss {
		if len(v) != 3 {
			continue
		}
		if strings.HasPrefix(v[2], "<![CDATA[") && strings.HasSuffix(v[2], "]]>") {
			v[2] = v[2][9 : len(v[2])-3]
		}
		if v[1] == "return_code" {
			returnCode = v[2]
		}
		if v[1] == "return_msg" {
			returnMsg = v[2]
		}
		if v[1] == "result_code" {
			resultCode = v[2]
		}
		if v[1] == "err_code_des" {
			errCodeDes = v[2]
		}
		if v[1] == "prepay_id" {
			prepayID = v[2]
		}
		if v[1] == "sign" {
			sign = v[2]
			continue
		}
		if v[2] == "" {
			continue
		}
		ss1 = append(ss1, v[1]+"="+v[2])
	}
	if returnCode != "SUCCESS" {
		return "", errors.New(returnMsg)
	}
	sort.Strings(ss1)
	s = strings.Join(ss1, "&")
	s += "&key=" + key
	if strings.ToUpper(MD5(s)) != sign {
		return "", errors.New("Invalid Sign")
	}
	if resultCode != "SUCCESS" {
		return "", errors.New(errCodeDes)
	}
	return prepayID, nil
}

func WxMiniMakePrepayID(outTradeNo, totalFee int64, openid, spbillCreateIP, body, notifyURL, appid, mchID, key string) (string, error) {
	nonceStr := strconv.FormatInt(RandomNumber(), 10)
	s := `
	<xml>
		<appid>%s</appid>
		<mch_id>%s</mch_id>
		<nonce_str>%s</nonce_str>
		<sign>%s</sign>
		<body>%s</body>
		<out_trade_no>%d</out_trade_no>
		<total_fee>%d</total_fee>
		<spbill_create_ip>%s</spbill_create_ip>
		<notify_url>%s</notify_url>
		<trade_type>JSAPI</trade_type>
		<openid>%s</openid>
	</xml>
	`
	s0 := fmt.Sprintf("appid=%s&body=%s&mch_id=%s&nonce_str=%s&notify_url=%s&openid=%s&out_trade_no=%d&spbill_create_ip=%s&total_fee=%d&trade_type=JSAPI&key=%s", appid, body, mchID, nonceStr, notifyURL, openid, outTradeNo, spbillCreateIP, totalFee, key)
	sign := strings.ToUpper(MD5(s0))
	s = fmt.Sprintf(s, appid, mchID, nonceStr, sign, body, outTradeNo, totalFee, spbillCreateIP, notifyURL, openid)

	c := &http.Client{
		Timeout: 10 * time.Second,
	}
	res, err := c.Post("https://api.mch.weixin.qq.com/pay/unifiedorder", "application/xml", bytes.NewReader([]byte(s)))
	if err != nil {
		return "", err
	}
	defer res.Body.Close()
	b, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return "", err
	}

	reg := regexp.MustCompile(`<(.+)>(.+)<`)
	ss := reg.FindAllStringSubmatch(strings.Replace(strings.Replace(string(b), "<xml>", "", -1), "</xml>", "", -1), -1)
	ss1 := make([]string, 0)
	var returnCode, returnMsg, resultCode, errCodeDes, prepayID string
	for _, v := range ss {
		if len(v) != 3 {
			continue
		}
		if strings.HasPrefix(v[2], "<![CDATA[") && strings.HasSuffix(v[2], "]]>") {
			v[2] = v[2][9 : len(v[2])-3]
		}
		if v[1] == "return_code" {
			returnCode = v[2]
		}
		if v[1] == "return_msg" {
			returnMsg = v[2]
		}
		if v[1] == "result_code" {
			resultCode = v[2]
		}
		if v[1] == "err_code_des" {
			errCodeDes = v[2]
		}
		if v[1] == "prepay_id" {
			prepayID = v[2]
		}
		if v[1] == "sign" {
			sign = v[2]
			continue
		}
		if v[2] == "" {
			continue
		}
		ss1 = append(ss1, v[1]+"="+v[2])
	}
	if returnCode != "SUCCESS" {
		return "", errors.New(returnMsg)
	}
	sort.Strings(ss1)
	s = strings.Join(ss1, "&")
	s += "&key=" + key
	if strings.ToUpper(MD5(s)) != sign {
		return "", errors.New("Invalid Sign")
	}
	if resultCode != "SUCCESS" {
		return "", errors.New(errCodeDes)
	}
	return prepayID, nil
}

type WxAppPayInfo struct {
	AppID      string
	PartnerID  string
	PrepayID   string
	PayPackage string
	NonceStr   string
	TimeStamp  string
	PaySign    string
}

func WxAppMakePayInfo(prepayid, appid, partnerid, key string) (*WxAppPayInfo, error) {
	ts := strconv.FormatInt(time.Now().Unix(), 10)
	nonceStr := strconv.FormatInt(RandomNumber(), 10)
	sign := strings.ToUpper(MD5(fmt.Sprintf("appid=%s&noncestr=%s&package=Sign=WXPay&partnerid=%s&prepayid=%s&timestamp=%s&key=%s", appid, nonceStr, partnerid, prepayid, ts, key)))

	o := &WxAppPayInfo{
		AppID:      appid,
		NonceStr:   nonceStr,
		PayPackage: "Sign=WXPay",
		PartnerID:  partnerid,
		PrepayID:   prepayid,
		TimeStamp:  ts,
		PaySign:    sign,
	}
	return o, nil
}

type WxMiniPayInfo struct {
	TimeStamp  string
	NonceStr   string
	PayPackage string
	SignType   string
	PaySign    string
}

func WxMiniMakePayInfo(prepayID, appId, key string) (*WxMiniPayInfo, error) {
	prepayID = "prepay_id=" + prepayID
	ts := strconv.FormatInt(time.Now().Unix(), 10)
	nonceStr := strconv.FormatInt(RandomNumber(), 10)
	sign := strings.ToUpper(MD5(fmt.Sprintf("appId=%s&nonceStr=%s&package=%s&signType=MD5&timeStamp=%s&key=%s", appId, nonceStr, prepayID, ts, key)))

	o := &WxMiniPayInfo{
		TimeStamp:  ts,
		NonceStr:   nonceStr,
		PayPackage: prepayID,
		SignType:   "MD5",
		PaySign:    sign,
	}
	return o, nil
}

func WxOrderCallbackFailBody(err error) string {
	fail := `
		<xml>
			<return_code><![CDATA[FAIL]]></return_code>
			<return_msg><![CDATA[%s]]></return_msg>
		</xml>
		`
	return fmt.Sprintf(fail, err.Error())
}
func WxOrderCallbackSuccessBody() string {
	success := `
		<xml>
			<return_code><![CDATA[SUCCESS]]></return_code>
			<return_msg><![CDATA[OK]]></return_msg>
		</xml>
		`
	return success
}

func WxOrderCallback(r *http.Request, key string) (int64, string, error) {
	b, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return 0, "", err
	}
	reg := regexp.MustCompile(`<(.+)>(.+)<`)
	ss := reg.FindAllStringSubmatch(strings.Replace(strings.Replace(string(b), "<xml>", "", -1), "</xml>", "", -1), -1)
	ss1 := make([]string, 0)
	var sign string
	var returnCode, returnMsg, resultCode, errCodeDes, outTradeNo, transactionID string
	for _, v := range ss {
		if len(v) != 3 {
			continue
		}
		if strings.HasPrefix(v[2], "<![CDATA[") && strings.HasSuffix(v[2], "]]>") {
			v[2] = v[2][9 : len(v[2])-3]
		}
		if v[1] == "return_code" {
			returnCode = v[2]
		}
		if v[1] == "return_msg" {
			returnMsg = v[2]
		}
		if v[1] == "result_code" {
			resultCode = v[2]
		}
		if v[1] == "err_code_des" {
			errCodeDes = v[2]
		}
		if v[1] == "out_trade_no" {
			outTradeNo = v[2]
		}
		if v[1] == "transaction_id" {
			transactionID = v[2]
		}
		if v[1] == "sign" {
			sign = v[2]
			continue
		}
		if v[2] == "" {
			continue
		}
		ss1 = append(ss1, v[1]+"="+v[2])
	}
	if returnCode != "SUCCESS" {
		return 0, "", errors.New(returnMsg)
	}
	sort.Strings(ss1)
	s := strings.Join(ss1, "&")
	s += "&key=" + key

	if strings.ToUpper(MD5(s)) != sign {
		return 0, "", errors.New("Invalid Sign")
	}
	if resultCode != "SUCCESS" {
		return 0, "", errors.New(errCodeDes)
	}
	id, err := strconv.ParseInt(outTradeNo, 10, 64)
	if err != nil {
		return 0, "", err
	}
	return id, transactionID, nil
}

func WxAppRefund(outTradeNo, outRefundNo, totalFee, refundFee int64, appid, mchID, key, certPEM, keyPEM string) (string, error) {
	cert, err := tls.X509KeyPair([]byte(certPEM), []byte(keyPEM))
	if err != nil {
		return "", err
	}
	tlsConfig := &tls.Config{
		Certificates: []tls.Certificate{cert},
	}
	c := &http.Client{
		Transport: &http.Transport{
			Proxy: http.ProxyFromEnvironment,
			DialContext: (&net.Dialer{
				Timeout:   3 * time.Second,
				KeepAlive: 30 * time.Second,
				DualStack: true,
			}).DialContext,
			MaxIdleConns:          100,
			IdleConnTimeout:       90 * time.Second,
			TLSHandshakeTimeout:   3 * time.Second,
			ExpectContinueTimeout: 1 * time.Second,
			TLSClientConfig:       tlsConfig,
		},
		Timeout: 10 * time.Second,
	}
	nonceStr := strconv.FormatInt(RandomNumber(), 10)
	s := `
	<xml>
		<appid>%s</appid>
		<mch_id>%s</mch_id>
		<nonce_str>%s</nonce_str>
		<sign>%s</sign>
		<out_trade_no>%d</out_trade_no>
		<out_refund_no>%d</out_refund_no>
		<total_fee>%d</total_fee>
		<refund_fee>%d</refund_fee>
	</xml>
	`
	s0 := fmt.Sprintf("appid=%s&mch_id=%s&nonce_str=%s&out_refund_no=%d&out_trade_no=%d&refund_fee=%d&total_fee=%d&key=%s", appid, mchID, nonceStr, outRefundNo, outTradeNo, refundFee, totalFee, key)
	sign := strings.ToUpper(MD5(s0))
	s = fmt.Sprintf(s, appid, mchID, nonceStr, sign, outTradeNo, outRefundNo, totalFee, refundFee)

	res, err := c.Post("https://api.mch.weixin.qq.com/secapi/pay/refund", "application/xml", bytes.NewReader([]byte(s)))
	if err != nil {
		return "", err
	}
	defer res.Body.Close()
	b, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return "", err
	}

	reg := regexp.MustCompile(`<(.+)>(.+)<`)
	ss := reg.FindAllStringSubmatch(strings.Replace(strings.Replace(string(b), "<xml>", "", -1), "</xml>", "", -1), -1)
	var returnCode, returnMsg, resultCode, errCodeDes, transactionId string
	for _, v := range ss {
		if len(v) != 3 {
			continue
		}
		if strings.HasPrefix(v[2], "<![CDATA[") && strings.HasSuffix(v[2], "]]>") {
			v[2] = v[2][9 : len(v[2])-3]
		}
		if v[1] == "return_code" {
			returnCode = v[2]
		}
		if v[1] == "return_msg" {
			returnMsg = v[2]
		}
		if v[1] == "result_code" {
			resultCode = v[2]
		}
		if v[1] == "err_code_des" {
			errCodeDes = v[2]
		}
		if v[1] == "transaction_id" {
			transactionId = v[2]
		}
	}
	if returnCode != "SUCCESS" {
		return "", errors.New(returnMsg)
	}
	if resultCode != "SUCCESS" {
		return "", errors.New(errCodeDes)
	}
	return transactionId, nil
}

func WxMiniRefund(outTradeNo, outRefundNo, totalFee, refundFee int64, appid, mchID, key, certPEM, keyPEM string) (string, error) {
	cert, err := tls.X509KeyPair([]byte(certPEM), []byte(keyPEM))
	if err != nil {
		return "", err
	}
	tlsConfig := &tls.Config{
		Certificates: []tls.Certificate{cert},
	}
	c := &http.Client{
		Transport: &http.Transport{
			Proxy: http.ProxyFromEnvironment,
			DialContext: (&net.Dialer{
				Timeout:   3 * time.Second,
				KeepAlive: 30 * time.Second,
				DualStack: true,
			}).DialContext,
			MaxIdleConns:          100,
			IdleConnTimeout:       90 * time.Second,
			TLSHandshakeTimeout:   3 * time.Second,
			ExpectContinueTimeout: 1 * time.Second,
			TLSClientConfig:       tlsConfig,
		},
		Timeout: 10 * time.Second,
	}
	nonceStr := strconv.FormatInt(RandomNumber(), 10)
	s := `
	<xml>
		<appid>%s</appid>
		<mch_id>%s</mch_id>
		<nonce_str>%s</nonce_str>
		<sign>%s</sign>
		<out_trade_no>%d</out_trade_no>
		<out_refund_no>%d</out_refund_no>
		<total_fee>%d</total_fee>
		<refund_fee>%d</refund_fee>
	</xml>
	`
	s0 := fmt.Sprintf("appid=%s&mch_id=%s&nonce_str=%s&out_refund_no=%d&out_trade_no=%d&refund_fee=%d&total_fee=%d&key=%s", appid, mchID, nonceStr, outRefundNo, outTradeNo, refundFee, totalFee, key)
	sign := strings.ToUpper(MD5(s0))
	s = fmt.Sprintf(s, appid, mchID, nonceStr, sign, outTradeNo, outRefundNo, totalFee, refundFee)

	res, err := c.Post("https://api.mch.weixin.qq.com/secapi/pay/refund", "application/xml", bytes.NewReader([]byte(s)))
	if err != nil {
		return "", err
	}
	defer res.Body.Close()
	b, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return "", err
	}

	reg := regexp.MustCompile(`<(.+)>(.+)<`)
	ss := reg.FindAllStringSubmatch(strings.Replace(strings.Replace(string(b), "<xml>", "", -1), "</xml>", "", -1), -1)
	var returnCode, returnMsg, resultCode, errCodeDes, transactionId string
	for _, v := range ss {
		if len(v) != 3 {
			continue
		}
		if strings.HasPrefix(v[2], "<![CDATA[") && strings.HasSuffix(v[2], "]]>") {
			v[2] = v[2][9 : len(v[2])-3]
		}
		if v[1] == "return_code" {
			returnCode = v[2]
		}
		if v[1] == "return_msg" {
			returnMsg = v[2]
		}
		if v[1] == "result_code" {
			resultCode = v[2]
		}
		if v[1] == "err_code_des" {
			errCodeDes = v[2]
		}
		if v[1] == "transaction_id" {
			transactionId = v[2]
		}
	}
	if returnCode != "SUCCESS" {
		return "", errors.New(returnMsg)
	}
	if resultCode != "SUCCESS" {
		return "", errors.New(errCodeDes)
	}
	return transactionId, nil
}

func WxAppTransfer(openid string, partnerTradeNo, amount int64, desc, mchAppid, mchid, key, certPEM, keyPEM string) (string, error) {
	cert, err := tls.X509KeyPair([]byte(certPEM), []byte(keyPEM))
	if err != nil {
		return "", err
	}
	tlsConfig := &tls.Config{
		Certificates: []tls.Certificate{cert},
	}
	c := &http.Client{
		Transport: &http.Transport{
			Proxy: http.ProxyFromEnvironment,
			DialContext: (&net.Dialer{
				Timeout:   3 * time.Second,
				KeepAlive: 30 * time.Second,
				DualStack: true,
			}).DialContext,
			MaxIdleConns:          100,
			IdleConnTimeout:       90 * time.Second,
			TLSHandshakeTimeout:   3 * time.Second,
			ExpectContinueTimeout: 1 * time.Second,
			TLSClientConfig:       tlsConfig,
		},
		Timeout: 10 * time.Second,
	}
	nonceStr := strconv.FormatInt(RandomNumber(), 10)
	s := `
	<xml>
		<mch_appid>%s</mch_appid>
		<mchid>%s</mchid>
		<nonce_str>%s</nonce_str>
		<sign>%s</sign>
		<partner_trade_no>%d</partner_trade_no>
		<openid>%s</openid>
		<check_name>NO_CHECK</check_name>
		<amount>%d</amount>
		<desc>%s</desc>
		<spbill_create_ip>8.8.8.8</spbill_create_ip>
	</xml>
	`
	s0 := fmt.Sprintf("amount=%d&check_name=NO_CHECK&desc=%s&mch_appid=%s&mchid=%s&nonce_str=%s&openid=%s&partner_trade_no=%d&spbill_create_ip=8.8.8.8&key=%s", amount, desc, mchAppid, mchid, nonceStr, openid, partnerTradeNo, key)
	sign := strings.ToUpper(MD5(s0))
	s = fmt.Sprintf(s, mchAppid, mchid, nonceStr, sign, partnerTradeNo, openid, amount, desc)

	res, err := c.Post("https://api.mch.weixin.qq.com/mmpaymkttransfers/promotion/transfers", "application/xml", bytes.NewReader([]byte(s)))
	if err != nil {
		return "", err
	}
	defer res.Body.Close()
	b, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return "", err
	}

	reg := regexp.MustCompile(`<(.+)>(.+)<`)
	ss := reg.FindAllStringSubmatch(strings.Replace(strings.Replace(string(b), "<xml>", "", -1), "</xml>", "", -1), -1)
	var returnCode, returnMsg, resultCode, errCodeDes, paymentNo string
	for _, v := range ss {
		if len(v) != 3 {
			continue
		}
		if strings.HasPrefix(v[2], "<![CDATA[") && strings.HasSuffix(v[2], "]]>") {
			v[2] = v[2][9 : len(v[2])-3]
		}
		if v[1] == "return_code" {
			returnCode = v[2]
		}
		if v[1] == "return_msg" {
			returnMsg = v[2]
		}
		if v[1] == "result_code" {
			resultCode = v[2]
		}
		if v[1] == "err_code_des" {
			errCodeDes = v[2]
		}
		if v[1] == "payment_no" {
			paymentNo = v[2]
		}
	}
	if returnCode != "SUCCESS" {
		return "", errors.New(returnMsg)
	}
	if resultCode != "SUCCESS" {
		return "", errors.New(errCodeDes)
	}
	return paymentNo, nil
}

func WxMiniTransfer(openid string, partnerTradeNo, amount int64, desc, mchAppid, mchid, key, certPEM, keyPEM string) (string, error) {
	cert, err := tls.X509KeyPair([]byte(certPEM), []byte(keyPEM))
	if err != nil {
		return "", err
	}
	tlsConfig := &tls.Config{
		Certificates: []tls.Certificate{cert},
	}
	c := &http.Client{
		Transport: &http.Transport{
			Proxy: http.ProxyFromEnvironment,
			DialContext: (&net.Dialer{
				Timeout:   3 * time.Second,
				KeepAlive: 30 * time.Second,
				DualStack: true,
			}).DialContext,
			MaxIdleConns:          100,
			IdleConnTimeout:       90 * time.Second,
			TLSHandshakeTimeout:   3 * time.Second,
			ExpectContinueTimeout: 1 * time.Second,
			TLSClientConfig:       tlsConfig,
		},
		Timeout: 10 * time.Second,
	}
	nonceStr := strconv.FormatInt(RandomNumber(), 10)
	s := `
	<xml>
		<mch_appid>%s</mch_appid>
		<mchid>%s</mchid>
		<nonce_str>%s</nonce_str>
		<sign>%s</sign>
		<partner_trade_no>%d</partner_trade_no>
		<openid>%s</openid>
		<check_name>NO_CHECK</check_name>
		<amount>%d</amount>
		<desc>%s</desc>
		<spbill_create_ip>8.8.8.8</spbill_create_ip>
	</xml>
	`
	s0 := fmt.Sprintf("amount=%d&check_name=NO_CHECK&desc=%s&mch_appid=%s&mchid=%s&nonce_str=%s&openid=%s&partner_trade_no=%d&spbill_create_ip=8.8.8.8&key=%s", amount, desc, mchAppid, mchid, nonceStr, openid, partnerTradeNo, key)
	sign := strings.ToUpper(MD5(s0))
	s = fmt.Sprintf(s, mchAppid, mchid, nonceStr, sign, partnerTradeNo, openid, amount, desc)

	res, err := c.Post("https://api.mch.weixin.qq.com/mmpaymkttransfers/promotion/transfers", "application/xml", bytes.NewReader([]byte(s)))
	if err != nil {
		return "", err
	}
	defer res.Body.Close()
	b, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return "", err
	}

	reg := regexp.MustCompile(`<(.+)>(.+)<`)
	ss := reg.FindAllStringSubmatch(strings.Replace(strings.Replace(string(b), "<xml>", "", -1), "</xml>", "", -1), -1)
	var returnCode, returnMsg, resultCode, errCodeDes, paymentNo string
	for _, v := range ss {
		if len(v) != 3 {
			continue
		}
		if strings.HasPrefix(v[2], "<![CDATA[") && strings.HasSuffix(v[2], "]]>") {
			v[2] = v[2][9 : len(v[2])-3]
		}
		if v[1] == "return_code" {
			returnCode = v[2]
		}
		if v[1] == "return_msg" {
			returnMsg = v[2]
		}
		if v[1] == "result_code" {
			resultCode = v[2]
		}
		if v[1] == "err_code_des" {
			errCodeDes = v[2]
		}
		if v[1] == "payment_no" {
			paymentNo = v[2]
		}
	}
	if returnCode != "SUCCESS" {
		return "", errors.New(returnMsg)
	}
	if resultCode != "SUCCESS" {
		return "", errors.New(errCodeDes)
	}
	return paymentNo, nil
}

func WxMPTransfer(openid string, partnerTradeNo, amount int64, desc, mchAppid, mchid, key, certPEM, keyPEM string) (string, error) {
	cert, err := tls.X509KeyPair([]byte(certPEM), []byte(keyPEM))
	if err != nil {
		return "", err
	}
	tlsConfig := &tls.Config{
		Certificates: []tls.Certificate{cert},
	}
	c := &http.Client{
		Transport: &http.Transport{
			Proxy: http.ProxyFromEnvironment,
			DialContext: (&net.Dialer{
				Timeout:   3 * time.Second,
				KeepAlive: 30 * time.Second,
				DualStack: true,
			}).DialContext,
			MaxIdleConns:          100,
			IdleConnTimeout:       90 * time.Second,
			TLSHandshakeTimeout:   3 * time.Second,
			ExpectContinueTimeout: 1 * time.Second,
			TLSClientConfig:       tlsConfig,
		},
		Timeout: 5 * time.Second,
	}
	nonceStr := strconv.FormatInt(RandomNumber(), 10)
	s := `
	<xml>
		<mch_appid>%s</mch_appid>
		<mchid>%s</mchid>
		<nonce_str>%s</nonce_str>
		<sign>%s</sign>
		<partner_trade_no>%d</partner_trade_no>
		<openid>%s</openid>
		<check_name>NO_CHECK</check_name>
		<amount>%d</amount>
		<desc>%s</desc>
		<spbill_create_ip>8.8.8.8</spbill_create_ip>
	</xml>
	`
	s0 := fmt.Sprintf("amount=%d&check_name=NO_CHECK&desc=%s&mch_appid=%s&mchid=%s&nonce_str=%s&openid=%s&partner_trade_no=%d&spbill_create_ip=8.8.8.8&key=%s", amount, desc, mchAppid, mchid, nonceStr, openid, partnerTradeNo, key)
	sign := strings.ToUpper(MD5(s0))
	s = fmt.Sprintf(s, mchAppid, mchid, nonceStr, sign, partnerTradeNo, openid, amount, desc)

	res, err := c.Post("https://api.mch.weixin.qq.com/mmpaymkttransfers/promotion/transfers", "application/xml", bytes.NewReader([]byte(s)))
	if err != nil {
		return "", err
	}
	defer res.Body.Close()
	b, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return "", err
	}

	reg := regexp.MustCompile(`<(.+)>(.+)<`)
	ss := reg.FindAllStringSubmatch(strings.Replace(strings.Replace(string(b), "<xml>", "", -1), "</xml>", "", -1), -1)
	var returnCode, returnMsg, resultCode, errCodeDes, paymentNo string
	for _, v := range ss {
		if len(v) != 3 {
			continue
		}
		if strings.HasPrefix(v[2], "<![CDATA[") && strings.HasSuffix(v[2], "]]>") {
			v[2] = v[2][9 : len(v[2])-3]
		}
		if v[1] == "return_code" {
			returnCode = v[2]
		}
		if v[1] == "return_msg" {
			returnMsg = v[2]
		}
		if v[1] == "result_code" {
			resultCode = v[2]
		}
		if v[1] == "err_code_des" {
			errCodeDes = v[2]
		}
		if v[1] == "payment_no" {
			paymentNo = v[2]
		}
	}
	if returnCode != "SUCCESS" {
		return "", errors.New(returnMsg)
	}
	if resultCode != "SUCCESS" {
		return "", errors.New(errCodeDes)
	}
	return paymentNo, nil
}

func WxMiniGetOpenIDByCode(code, appid, secret string) (string, error) {
	client := &http.Client{
		Timeout: 10 * time.Second,
	}
	r, err := client.Get(fmt.Sprintf("https://api.weixin.qq.com/sns/jscode2session?appid=%s&secret=%s&js_code=%s&grant_type=authorization_code", appid, secret, code))
	if err != nil {
		return "", err
	}
	defer r.Body.Close()
	j, err := simplejson.NewFromReader(r.Body)
	if err != nil {
		return "", err
	}
	s, err := j.Get("errmsg").String()
	if err == nil {
		return "", errors.New(s)
	}
	s, err = j.Get("openid").String()
	if err != nil {
		return "", err
	}
	return s, nil
}

func WxMPOAuthGetOpenIDAndAccessToken(code, appid, secret string) (string, string, error) {
	client := &http.Client{
		Timeout: 10 * time.Second,
	}
	r, err := client.Get(fmt.Sprintf("https://api.weixin.qq.com/sns/oauth2/access_token?appid=%s&secret=%s&code=%s&grant_type=authorization_code", appid, secret, code))
	if err != nil {
		return "", "", err
	}
	defer r.Body.Close()
	j, err := simplejson.NewFromReader(r.Body)
	if err != nil {
		return "", "", err
	}
	s, err := j.Get("errmsg").String()
	if err == nil {
		return "", "", errors.New(s)
	}
	s, err = j.Get("openid").String()
	if err != nil {
		return "", "", err
	}
	s1, err := j.Get("access_token").String()
	if err != nil {
		return "", "", err
	}
	return s, s1, nil
}

type WxMPUser struct {
	OpenID     string `json: "openid"`
	NickName   string `json: "nickname"`
	Sex        int64  `json: "sex"`
	Province   string `json: "province"`
	City       string `json: "city"`
	Country    string `json: "country"`
	HeadImgUrl string `json: "headimgurl"`
}

func WxMPGetUserInfo(openid, accessToken string) (*WxMPUser, error) {
	client := &http.Client{
		Timeout: 10 * time.Second,
	}
	r, err := client.Get(fmt.Sprintf("https://api.weixin.qq.com/sns/userinfo?access_token=%s&openid=%s&lang=zh_CN", accessToken, openid))
	if err != nil {
		return nil, err
	}
	defer r.Body.Close()
	j, err := simplejson.NewFromReader(r.Body)
	if err != nil {
		return nil, err
	}
	s, err := j.Get("errmsg").String()
	if err == nil {
		return nil, errors.New(s)
	}
	b, err := j.MarshalJSON()
	if err != nil {
		return nil, err
	}
	o := &WxMPUser{}
	if err := json.Unmarshal(b, o); err != nil {
		return nil, err
	}
	return o, nil
}
