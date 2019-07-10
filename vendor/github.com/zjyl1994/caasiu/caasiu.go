package caasiu

import (
	"io/ioutil"
	"net/http"
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
	defer r.Body.Close()
	s, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return nil, err
	}
	caasiu.jsonBody = NewJSON(s)
	return &caasiu, nil
}

func (c *Caasiu) Body() *JSON {
	return c.jsonBody
}

func (c *Caasiu) QueryString() *QueryString {
	return c.queryString
}

func (c *Caasiu) Valid(rules Rules) (bool, []string) {
	var errMsg []string
	if len(rules.QueryString) > 0 {
		errMsg = append(errMsg, c.queryString.Valid(rules.QueryString)...)
	}
	if len(rules.Body) > 0 {
		errMsg = append(errMsg, c.jsonBody.Valid(rules.Body)...)
	}
	if len(errMsg) > 0 {
		return false, errMsg
	} else {
		return true, nil
	}
}
