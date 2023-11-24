package x

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
)

func ReadJSON(r *http.Request, o interface{}) error {
	d, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return err
	}
	if err = json.Unmarshal(d, o); err != nil {
		return err
	}
	return nil
}

func JSON(w http.ResponseWriter, v interface{}) {
	d, err := json.Marshal(v)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	w.Header().Set("content-type", "application/json")
	w.Write(d)
}
