package uroncha

import (
	simplejson "github.com/bitly/go-simplejson"
	"github.com/gin-gonic/gin"
	"github.com/zjyl1994/caasiu"
)

type H = gin.H
type Context = gin.Context
type Rule = caasiu.Rule
type Rules = caasiu.Rules

type HandleFunc func(*Context, Datas) (interface{}, Error)

type Datas struct {
	QueryString map[string]string
	Body        *simplejson.Json
}

type Error struct {
	Success bool   `json:"success"`
	Code    int    `json:"code"`
	Message string `json:"message"`
}

var NoError = Error{
	Success: true,
	Code:    0,
	Message: "",
}
