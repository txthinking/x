package ant

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"path/filepath"
)

// TODO combine

// RFC 2388
func MultipartFormDataFromFile(params, files map[string][]string, boundary string) (ior io.Reader, err error) {
	var bs []byte
	var bf *bytes.Buffer = &bytes.Buffer{}

	// prepare common value
	var name, value string
	var values []string
	for name, values = range params {
		for _, value = range values {
			bf.WriteString(fmt.Sprintf("--%s\r\n", boundary))
			bf.WriteString(fmt.Sprintf("Content-Disposition: form-data; name=\"%s\"\r\n\r\n", name))
			bf.WriteString(fmt.Sprintf("%s\r\n", value))
		}
	}

	for name, values = range files {
		for _, value = range values {
			bf.WriteString(fmt.Sprintf("--%s\r\n", boundary))
			bf.WriteString(fmt.Sprintf("Content-Disposition: form-data; name=\"%s\"; filename=\"%s\"\r\n", name, filepath.Base(value)))
			bf.WriteString(fmt.Sprintf("Content-Type: application/octet-stream\r\n\r\n"))
			bs, err = ioutil.ReadFile(value)
			if err != nil {
				return
			}
			bf.Write(bs)
			bf.WriteString("\r\n")
		}
	}
	bf.WriteString(fmt.Sprintf("--%s--\r\n", boundary))
	ior = bf
	return
}

// RFC 2388
func MultipartFormDataFromReader(params map[string][]string, files map[string][]io.Reader, boundary string) (ior io.Reader, err error) {
	var bs []byte
	var bf *bytes.Buffer = &bytes.Buffer{}

	// prepare common value
	var name, value string
	var values []string
	for name, values = range params {
		for _, value = range values {
			bf.WriteString(fmt.Sprintf("--%s\r\n", boundary))
			bf.WriteString(fmt.Sprintf("Content-Disposition: form-data; name=\"%s\"\r\n\r\n", name))
			bf.WriteString(fmt.Sprintf("%s\r\n", value))
		}
	}

	var rs []io.Reader
	var r io.Reader
	for name, rs = range files {
		for _, r = range rs {
			bf.WriteString(fmt.Sprintf("--%s\r\n", boundary))
			bf.WriteString(fmt.Sprintf("Content-Disposition: form-data; name=\"%s\"; filename=\"%s\"\r\n", name, "-"))
			bf.WriteString(fmt.Sprintf("Content-Type: application/octet-stream\r\n\r\n"))
			bs, err = ioutil.ReadAll(r)
			if err != nil {
				return
			}
			bf.Write(bs)
			bf.WriteString("\r\n")
		}
	}
	bf.WriteString(fmt.Sprintf("--%s--\r\n", boundary))
	ior = bf
	return
}
