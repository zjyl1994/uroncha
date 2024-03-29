package uroncha

import (
	"github.com/zjyl1994/caasiu"
)

type H map[string]interface{}
type Rule = caasiu.Rule
type Rules = caasiu.Rules

type HandleFunc func(*Context) (interface{}, Error)

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

type DownloadFile struct {
	FilePath    string
	ContentType string
	FileName    string
}

func NewError(err error) Error {
	if err == nil {
		return NoError
	} else {
		return Error{
			Success: false,
			Code:    -1,
			Message: err.Error(),
		}
	}
}
