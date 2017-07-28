package ant

import (
	"bytes"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	sj "github.com/bitly/go-simplejson"
	"github.com/fatih/structs"
)

type Slack struct {
	URL      string
	Channel  string
	UserName string
	IconURL  string
}

// Send data to slack channel
// data must be one of map[string]string, map[string]interface{}, string, []string, struct
// cui is channel like #general, username, icon_url
func (sl *Slack) Send(data interface{}, cui ...string) error {
	if structs.IsStruct(data) {
		data = structs.Map(data)
	}
	var s string
	switch data.(type) {
	case string:
		s = data.(string)
	case []string:
		for i, v := range data.([]string) {
			s += fmt.Sprintf("%v: %v\n", i, v)
		}
	case map[string]string:
		for k, v := range data.(map[string]string) {
			s += fmt.Sprintf("%v: %v\n", k, v)
		}
	case map[string]interface{}:
		for k, v := range data.(map[string]interface{}) {
			s += fmt.Sprintf("%v: %v\n", k, v)
		}
	default:
		return errors.New("Unsupport data")
	}

	if len(cui) == 1 {
		sl.Channel = cui[0]
	}
	if len(cui) == 2 {
		sl.Channel = cui[0]
		sl.UserName = cui[1]
	}
	if len(cui) == 3 {
		sl.Channel = cui[0]
		sl.UserName = cui[1]
		sl.IconURL = cui[2]
	}
	if sl.Channel == "" {
		sl.Channel = "#general"
	}
	if sl.UserName == "" {
		sl.UserName = "Bot"
	}

	j := sj.New()
	j.Set("text", "```"+s+"```")
	j.Set("channel", sl.Channel)
	j.Set("username", sl.UserName)
	if sl.IconURL != "" {
		j.Set("icon_url", sl.IconURL)
	}
	d, err := j.Encode()
	if err != nil {
		return err
	}

	client := &http.Client{
		Timeout: 10 * time.Second,
	}
	r, err := client.Post(sl.URL, "application/json", bytes.NewReader(d))
	if err != nil {
		return err
	}
	defer r.Body.Close()
	d, err = ioutil.ReadAll(r.Body)
	if err != nil {
		return err
	}
	if r.StatusCode != 200 {
		return errors.New(r.Status + ": " + string(d))
	}
	return nil
}
