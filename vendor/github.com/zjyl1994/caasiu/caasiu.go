package caasiu

import (
	"io/ioutil"
	"net/http"

	simplejson "github.com/bitly/go-simplejson"
)

type Caasiu struct {
	req         *http.Request
	jsonBody    *JSON
	queryString *QueryString
}

func New(r *http.Request) (*Caasiu, error) {
	var caasiu Caasiu
	caasiu.req = r
	caasiu.queryString = NewQueryString(r.URL.Query())
	contentType := r.Header.Get("Content-type")
	if contentType == "application/json" {
		defer r.Body.Close()
		s, err := ioutil.ReadAll(r.Body)
		if err != nil {
			return nil, err
		}
		caasiu.jsonBody = NewJSON(s)
	} else {
		caasiu.jsonBody = nil
	}
	return &caasiu, nil
}

func (c *Caasiu) QueryStringData() map[string]string {
	if c.queryString == nil {
		return nil
	} else {
		return c.queryString.Data()
	}
}

func (c *Caasiu) JsonBodyData() *simplejson.Json {
	if c.jsonBody == nil {
		return nil
	} else {
		return c.jsonBody.Data()
	}
}

func (c *Caasiu) Valid(rules Rules) (bool, []string) {
	var errMsg []string
	if len(rules.QueryString) > 0 {
		errMsg = append(errMsg, c.queryString.Valid(rules.QueryString)...)
	}
	if len(rules.Body) > 0 && c.jsonBody != nil {
		errMsg = append(errMsg, c.jsonBody.Valid(rules.Body)...)
	}
	if len(errMsg) > 0 {
		return false, errMsg
	} else {
		return true, nil
	}
}
