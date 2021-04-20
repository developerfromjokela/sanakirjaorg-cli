package networking

import (
	"encoding/json"
	"net/http"
	"net/url"
	"time"
)

var myClient = &http.Client{Timeout: 10 * time.Second}

func GetJson(urlStr string, target interface{}) error {
	var request = &http.Request{}
	uri, err2 := url.ParseRequestURI(urlStr)
	if err2 != nil {
		return err2
	}
	request.URL = uri
	request.Method = "GET"
	headerMap := map[string][]string{}
	headerMap["User-Agent"] = []string{"Mozilla/5.0 (sanakirja-cmd)"}
	request.Header = headerMap
	r, err := myClient.Do(request)
	if err != nil {
		return err
	}
	defer r.Body.Close()
	return json.NewDecoder(r.Body).Decode(target)
}
