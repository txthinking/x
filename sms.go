package x

import (
	"bytes"
	"encoding/base64"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	simplejson "github.com/bitly/go-simplejson"
)

type TencentSMS struct {
	AppID  string
	AppKey string
}

func (t *TencentSMS) Send(message, mobile string) error {
	c := &http.Client{
		Timeout: 10 * time.Second,
	}
	rd := strconv.FormatInt(RandomNumber(), 10)
	url := fmt.Sprintf("https://yun.tim.qq.com/v5/tlssmssvr/sendsms?sdkappid=%s&random=%s", t.AppID, rd)
	body := `
	{
		"ext": "",
		"extend": "",
		"msg": "%s",
		"sig": "%s",
		"tel": {
			"mobile": "%s",
			"nationcode": "86"
		},
		"time": %s,
		"type": 0
	}
	`
	ts := strconv.FormatInt(time.Now().Unix(), 10)
	sig := SHA256String(fmt.Sprintf("appkey=%s&random=%s&time=%s&mobile=%s", t.AppKey, rd, ts, mobile))
	body = fmt.Sprintf(body, message, sig, mobile, ts)
	rq, err := http.NewRequest("POST", url, strings.NewReader(body))
	if err != nil {
		return err
	}
	res, err := c.Do(rq)
	if err != nil {
		return err
	}
	defer res.Body.Close()
	if res.StatusCode != 200 {
		return errors.New(res.Status)
	}
	j, err := simplejson.NewFromReader(res.Body)
	if err != nil {
		return err
	}
	i, err := j.Get("result").Int()
	if err != nil {
		return err
	}
	if i != 0 {
		s, err := j.Get("errmsg").String()
		if err != nil {
			return err
		}
		return errors.New(s)
	}
	return nil
}

type YunpianSMS struct {
	ApiKey string
}

func (y *YunpianSMS) Send(message, mobile string) error {
	client := &http.Client{
		Timeout: 10 * time.Second,
	}
	v := url.Values{}
	v.Add("apikey", y.ApiKey)
	v.Add("mobile", mobile)
	v.Add("text", message)
	r, err := client.PostForm("https://sms.yunpian.com/v2/sms/single_send.json", v)
	if err != nil {
		return err
	}
	defer r.Body.Close()
	j, err := simplejson.NewFromReader(r.Body)
	if err != nil {
		return err
	}
	i, err := j.Get("code").Int()
	if err != nil {
		return err
	}
	if i == 0 {
		return nil
	}
	s, err := j.Get("msg").String()
	if err != nil {
		return err
	}
	return errors.New(s)
}

type JPush struct {
	AppKey       string
	MasterSecret string
}

// If ids length equal 0 means push to all,
// extras must can be encoded into json object.
func (jp *JPush) Notification(ids []string, alert string, extras interface{}, production bool) error {
	body := `
{
    "platform": "all",
    "audience": {
    },
    "notification": {
		"android": {
			"alert": "",
			"extras": {
			}
		},
		"ios": {
			"alert": "",
			"content-available": true,
			"extras": {
			}
		}
    },
    "options": {
        "apns_production": false
    }
}
	`
	j, err := simplejson.NewJson([]byte(body))
	if err != nil {
		return err
	}
	j.SetPath([]string{"audience", "registration_id"}, ids)
	if len(ids) == 0 {
		j.SetPath([]string{"audience"}, "all")
	}
	j.SetPath([]string{"notification", "android", "alert"}, alert)
	j.SetPath([]string{"notification", "ios", "alert"}, alert)
	j.SetPath([]string{"notification", "android", "extras"}, extras)
	j.SetPath([]string{"notification", "ios", "extras"}, extras)
	j.SetPath([]string{"options", "apns_production"}, production)
	b, err := j.MarshalJSON()
	if err != nil {
		return err
	}

	rq, err := http.NewRequest("POST", "https://api.jpush.cn/v3/push", bytes.NewReader(b))
	if err != nil {
		return err
	}
	rq.Header.Set("Authorization", "Basic "+base64.StdEncoding.EncodeToString([]byte(jp.AppKey+":"+jp.MasterSecret)))
	client := &http.Client{
		Timeout: 10 * time.Second,
	}
	r, err := client.Do(rq)
	if err != nil {
		return err
	}
	defer r.Body.Close()
	if r.StatusCode == 200 {
		return nil
	}

	j, err = simplejson.NewFromReader(r.Body)
	if err != nil {
		return err
	}
	s, err := j.GetPath("error", "message").String()
	if err != nil {
		return err
	}
	return errors.New(s)
}

// If ids length equal 0 means push to all,
// extras must can be encoded into json object.
func (jp *JPush) Message(ids []string, extras interface{}) error {
	body := `
{
    "platform": "all",
    "audience": {
    },
    "message": {
        "msg_content": "",
        "extras": {
        }
    }
}
	`
	j, err := simplejson.NewJson([]byte(body))
	if err != nil {
		return err
	}
	j.SetPath([]string{"audience", "registration_id"}, ids)
	if len(ids) == 0 {
		j.SetPath([]string{"audience"}, "all")
	}
	j.SetPath([]string{"message", "extras"}, extras)
	b, err := j.MarshalJSON()
	if err != nil {
		return err
	}

	rq, err := http.NewRequest("POST", "https://api.jpush.cn/v3/push", bytes.NewReader(b))
	if err != nil {
		return err
	}
	rq.Header.Set("Authorization", "Basic "+base64.StdEncoding.EncodeToString([]byte(jp.AppKey+":"+jp.MasterSecret)))
	client := &http.Client{
		Timeout: 10 * time.Second,
	}
	r, err := client.Do(rq)
	if err != nil {
		return err
	}
	defer r.Body.Close()
	if r.StatusCode == 200 {
		return nil
	}

	j, err = simplejson.NewFromReader(r.Body)
	if err != nil {
		return err
	}
	s, err := j.GetPath("error", "message").String()
	if err != nil {
		return err
	}
	return errors.New(s)
}

type TwilioSMS struct {
	Account string
	Token   string
	From    string
}

func (t *TwilioSMS) Send(message, mobile string) error {
	v := url.Values{}
	v.Set("From", t.From)
	v.Set("To", mobile)
	v.Set("Body", message)
	rq, err := http.NewRequest("POST", "https://api.twilio.com/2010-04-01/Accounts/"+t.Account+"/Messages.json", bytes.NewReader([]byte(v.Encode())))
	if err != nil {
		return err
	}
	rq.Header.Set("Authorization", "Basic "+base64.StdEncoding.EncodeToString([]byte(t.Account+":"+t.Token)))
	rq.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	client := &http.Client{
		Timeout: 10 * time.Second,
	}
	r, err := client.Do(rq)
	if err != nil {
		return err
	}
	defer r.Body.Close()
	if r.StatusCode == 201 {
		return nil
	}
	j, err := simplejson.NewFromReader(r.Body)
	if err != nil {
		return err
	}
	s, err := j.GetPath("message").String()
	if err != nil {
		return err
	}
	return errors.New(s)
}
