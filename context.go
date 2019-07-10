package uroncha

import (
	"net/http"

	simplejson "github.com/bitly/go-simplejson"
)

type Context struct {
	Req         *http.Request
	Res         http.ResponseWriter
	QueryString map[string]string
	PathParam   map[string]string
	Body        *simplejson.Json
}
