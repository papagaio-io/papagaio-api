package test

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
)

func parseBody(resp *http.Response, dto interface{}) {
	data, _ := ioutil.ReadAll(resp.Body)
	json.Unmarshal([]byte(string(data)), dto)
}
