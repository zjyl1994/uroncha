package uroncha

import (
	"io"
	"net/http"
	"os"

	simplejson "github.com/bitly/go-simplejson"
)

type Context struct {
	Req         *http.Request
	Res         http.ResponseWriter
	QueryString map[string]string
	PathParam   map[string]string
	Body        *simplejson.Json
}

func (c *Context) SaveFile(fieldName, savePath string) error {
	file, _, err := c.Req.FormFile(fieldName)
	if err != nil {
		return err
	}
	defer file.Close()
	saveFile, err := os.Create(savePath)
	if err != nil {
		return err
	}
	defer saveFile.Close()
	_, err = io.Copy(saveFile, file)
	return err
}
