package ant

import (
	"errors"
	"fmt"

	"github.com/fatih/structs"
)

// ToString convert data to string, data must be one of map[string]string, map[string]interface{}, string, []string, struct
func ToString(data interface{}) (string, error) {
	if structs.IsStruct(data) {
		data = structs.Map(data)
	}
	var s string
	switch data.(type) {
	case string:
		s = data.(string)
	case []string:
		for _, v := range data.([]string) {
			s += fmt.Sprintf("%v\n", v)
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
		return "", errors.New("Unsupport data")
	}
	return s, nil
}
