package main

import (
	"time"

	"github.com/zjyl1994/uroncha"
	"github.com/zjyl1994/uroncha/utils"
)

var logger = uroncha.Logger

func main() {
	logger.Println("Uroncha Example")
	uroncha.Handle("GET", "/ping", uroncha.Rules{}, pingHandler)
	uroncha.Handle("GET", "/get", uroncha.Rules{
		QueryString: uroncha.Rule{
			"type": []string{"required", "in:1,2,3,4,5"},
			"ts":   []string{"unix_timestamp"},
		},
	}, getHandler)
	uroncha.Handle("POST", "/post", uroncha.Rules{
		Body: uroncha.Rule{
			"data1":       []string{"required"},
			"data2.data3": []string{"required", "integer"},
			"data4":       []string{"string", "in:d1,d2,d3"},
		},
	}, postHandler)
	uroncha.Handle("GET", "/sleep", uroncha.Rules{
		QueryString: uroncha.Rule{
			"millisecond": []string{"required", "integer"},
		},
	}, sleepHandler)
	uroncha.Handle("GET", "/panic", uroncha.Rules{}, panicHandler)
	uroncha.Run()
}

//EXAMPLE: GET http://127.0.0.1:8080/ping
func pingHandler(c *uroncha.Context) (interface{}, uroncha.Error) {
	return uroncha.H{
		"name":  "URONCHA",
		"debug": uroncha.IsDebug(),
	}, uroncha.NoError
}

//EXAMPLE: GET http://127.0.0.1:8080/get?type=2&ts=1562752722
func getHandler(c *uroncha.Context) (interface{}, uroncha.Error) {
	return uroncha.H{
		"type": c.QueryString["type"],
		"ts":   c.QueryString["ts"],
	}, uroncha.NoError
}

/*
EXAMPLE: POST http://127.0.0.1:8080/post
{
	"data1":"asdfghjkl",
	"data2":{
		"data3":2
	},
	"data4":"d3"
}
*/
func postHandler(c *uroncha.Context) (interface{}, uroncha.Error) {
	body := c.Body
	data1 := body.Get("data1").MustString()
	data3 := body.Get("data2").Get("data3").MustInt()
	data4 := body.Get("data4").MustString()
	return uroncha.H{
		"data1": data1,
		"data3": data3,
		"data4": data4,
	}, uroncha.NoError
}

//EXAMPLE: GET http://127.0.0.1:8080/sleep?millisecond=100
func sleepHandler(c *uroncha.Context) (interface{}, uroncha.Error) {
	time.Sleep(time.Duration(utils.Str2Int(c.QueryString["millisecond"])) * time.Millisecond)
	return uroncha.H{
		"name":  "URONCHA",
		"debug": uroncha.IsDebug(),
	}, uroncha.NoError
}

func panicHandler(c *uroncha.Context) (interface{}, uroncha.Error) {
	panic("don't panic!")
}
