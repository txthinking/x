package x

import (
	"testing"
)

func TestMultipartFormDataFromFile(t *testing.T) {
	r, err := MultipartFormDataFromFile(
		map[string][]string{
			"a": {"哈啊好"},
		},
		map[string][]string{
			"f": {"/etc/hosts"},
		},
		"aasfsafsfsfsfsdf",
	)
	if err != nil {
		t.Fatal(r, err)
	}
	t.Log(r)
}
