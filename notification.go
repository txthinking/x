package x

import (
	"bytes"
	"errors"
	"fmt"
	"net/http"
	"time"

	simplejson "github.com/bitly/go-simplejson"
)

type Chat struct {
	URL string
}

// Send data to hangout chat
// data must be one of map[string]string, map[string]interface{}, string, []string, struct
func (c *Chat) SendText(data interface{}) error {
	s, err := ToString(data)
	if err != nil {
		return err
	}

	client := &http.Client{
		Timeout: 10 * time.Second,
	}
	j := simplejson.New()
	if _, ok := data.(string); ok {
		j.Set("text", s)
	} else {
		j.Set("text", "```\n"+s+"\n```")
	}
	b, err := j.MarshalJSON()
	if err != nil {
		return err
	}
	rq, err := http.NewRequest("POST", c.URL, bytes.NewReader(b))
	if err != nil {
		return err
	}
	r, err := client.Do(rq)
	if err != nil {
		return err
	}
	defer r.Body.Close()
	if r.StatusCode != 200 {
		return errors.New(r.Status)
	}
	return nil
}

func (c *Chat) SendCard(img, text, link string) error {
	client := &http.Client{
		Timeout: 10 * time.Second,
	}
	s := `
{
  "cards": [
    {
      "sections": [
        {
          "widgets": [
            {
              "image": { "imageUrl": "%s" }
            },
            {
              "buttons": [
                {
                  "textButton": {
                    "text": "%s",
                    "onClick": {
                      "openLink": {
                        "url": "%s"
                      }
                    }
                  }
                }
              ]
            }
          ]
        }
      ]
    }
  ]
}
`
	s = fmt.Sprintf(s, img, text, link)
	rq, err := http.NewRequest("POST", c.URL, bytes.NewBufferString(s))
	if err != nil {
		return err
	}
	r, err := client.Do(rq)
	if err != nil {
		return err
	}
	defer r.Body.Close()
	if r.StatusCode != 200 {
		return errors.New(r.Status)
	}
	return nil
}
