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
	route := &Route{HandlerFuncs: handlerFuncs, AutoSlash: true}
	route.RouteGroup.route = route
	return route
}

func (r *Route) AddBefore(hf HandlerFunc) {
	//append(r.Before, hf)
}

func (r *Route) AddGroup(path string) *RouteGroup {
	return &RouteGroup{path, r.route}
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
	for _, v := range r.Before {
		v(context)
	}

	v, ok := r.HandlerFuncs[re.URL.Path]
	if ok != true {
		r.redirectTo404(context)
	} else {
		v(context)
	}

}

type RouteGroup struct {
	groupPath string
	route     *Route
}

func (g *RouteGroup) Handle(path string, hf HandlerFunc) {
	var newPath string
	if g.route.AutoSlash == false {
		g.route.AddRoute(g.groupPath+path, hf)
	} else {
		if !strings.HasPrefix(path, "/") {
			newPath = g.groupPath + "/" + path
		} else {
			newPath = g.groupPath + path
		}
		if strings.HasSuffix(newPath, "/") {
			newPath = newPath[:len(newPath)-2]
		}
		g.route.AddRoute(newPath, hf)
	}

}
