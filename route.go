package ggweb

import (
	"log"
	"net/http"
	"strings"
)

type Out map[string]interface{}

type HandlerFunc func(*Context)

type Route struct {
	Before       []HandlerFunc
	HandlerFuncs map[string]HandlerFunc
	After        map[string]HandlerFunc
	AutoSlash    bool
	RouteGroup
}

func NewRoute() *Route {
	handlerFuncs := make(map[string]HandlerFunc)
	route := &Route{HandlerFuncs: handlerFuncs}
	route.RouteGroup.Route = route
	return route
}

func (r *Route) AddBefore(hf HandlerFunc) {
	//append(r.Before, hf)
}

func (r *Route) AddGroup(path string) *RouteGroup {
	return &RouteGroup{path, r.Route}
}

func (r *Route) AddRoute(path string, hf HandlerFunc) {
	if r.HandlerFuncs[path] != nil {
		log.Println("已存在")
		return
	}

	r.HandlerFuncs[path] = hf
	log.Println("-----------------------", r.HandlerFuncs)
}

func (r *Route) redirectTo404(c *Context) {
	type data struct {
		Content string
	}
	c.JSON(200, data{"404"})
}

func (r *Route) ServeHTTP(rw http.ResponseWriter, re *http.Request) {
	context := NewContext(rw, re)
	v, ok := r.HandlerFuncs[re.URL.Path]
	if ok != true {
		r.redirectTo404(context)
	} else {
		v(context)
	}

}

type RouteGroup struct {
	AbsolutePath string
	Route        *Route
}

func (g *RouteGroup) Handle(path string, hf HandlerFunc) {
	if g.Route.AutoSlash == false {
		g.Route.AddRoute(g.AbsolutePath+path, hf)
	} else {
		if !strings.HasPrefix(path, "/") {
			g.AbsolutePath = g.AbsolutePath + "/" + path
		}
		if strings.HasSuffix(path, "/") {
			g.AbsolutePath = g.AbsolutePath + path
			g.AbsolutePath = g.AbsolutePath[:len(g.AbsolutePath)-2]
			log.Println(g.AbsolutePath)
		}
		g.Route.AddRoute(g.AbsolutePath, hf)
	}

}
