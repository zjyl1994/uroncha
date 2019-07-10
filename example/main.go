package main

import "github.com/zjyl1994/uroncha"

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
	uroncha.Run()
}

//EXAMPLE: GET http://127.0.0.1:8000/ping
func pingHandler(c *uroncha.Context, params uroncha.Datas) (interface{}, uroncha.Error) {
	return uroncha.H{
		"name":  "URONCHA",
		"debug": uroncha.IsDebug(),
	}, uroncha.NoError
}

//EXAMPLE: GET http://127.0.0.1:8000/get?type=2&ts=1562752722
func getHandler(c *uroncha.Context, params uroncha.Datas) (interface{}, uroncha.Error) {
	return uroncha.H{
		"type": params.QueryString["type"],
		"ts":   params.QueryString["ts"],
	}, uroncha.NoError
}

/*
EXAMPLE: POST http://127.0.0.1:8000/post
{
	"data1":"asdfghjkl",
	"data2":{
		"data3":2
	},
	"data4":"d3"
}
*/
func postHandler(c *uroncha.Context, params uroncha.Datas) (interface{}, uroncha.Error) {
	body := params.Body
	data1 := body.Get("data1").MustString()
	data3 := body.Get("data2").Get("data3").MustInt()
	data4 := body.Get("data4").MustString()
	return uroncha.H{
		"data1": data1,
		"data3": data3,
		"data4": data4,
	}, uroncha.NoError
}
