package ant

import (
	"errors"
	"fmt"
	"net/http"
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
