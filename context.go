package ggweb

import (
	"encoding/json"
	"encoding/xml"
	"log"
	"net/http"
	"net/url"
)

type IContext interface {
	Json()
	String()
}

type Context struct {
	ResponseWriter http.ResponseWriter
	Request        *http.Request
	POST           map[string]string
	GET            map[string]string
	ALL            map[string]string
}

func NewContext(w http.ResponseWriter, r *http.Request) *Context {
	postForm := make(map[string]string)
	getForm := make(map[string]string)
	allForm := make(map[string]string)
	r.ParseForm()
	if r.Method == "GET" {
		values, _ := url.ParseQuery(r.URL.RawQuery)
		for k, v := range values {
			getForm[k] = v[0]
		}
	} else if r.Method == "POST" {
		values := r.PostForm
		for k, v := range values {
			postForm[k] = v[0]
		}
	}
	values := r.Form
	for k, v := range values {
		allForm[k] = v[0]
	}

	return &Context{ResponseWriter: w, Request: r, POST: postForm, GET: getForm, ALL: allForm}
}

func (ctx *Context) JSON(code int, i interface{}) {
	ctx.ResponseWriter.WriteHeader(code)
	ctx.ResponseWriter.Header().Set("Content-Type", "application/json; charset=utf-8")
	var res []byte
	switch i.(type) {
	default:
		var err error
		res, err = json.Marshal(&i)
		if err != nil {
			log.Println(err)
			return
		}
	case string:
		res = []byte(i.(string))
	}
	ctx.ResponseWriter.Write(res)
}

func (ctx *Context) STRING(code int, content string) {
	ctx.ResponseWriter.WriteHeader(code)
	ctx.ResponseWriter.Header().Set("Content-Type", "application/text; charset=utf-8")
	ctx.ResponseWriter.Write([]byte(content))
}

func (ctx *Context) XML(code int, i interface{}) {
	ctx.ResponseWriter.WriteHeader(code)
	ctx.ResponseWriter.Header().Set("Content-Type", "application/xml; charset=utf-8")
	var res []byte
	switch i.(type) {
	default:
		var err error
		res, err = xml.Marshal(&i)
		if err != nil {
			log.Println(err)
			return
		}
	case string:
		res = []byte(i.(string))
	}
	ctx.ResponseWriter.Write(res)
}
