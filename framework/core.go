package framework

import (
	"log"
	"net/http"
	"strings"
)

type Core struct {
	//router map[string]ControllerHandler
	router      map[string]*Tree
	middlewares []ControllerHandler
}

func NewCore() *Core {
	/*
		return &Core{
			router: map[string]ControllerHandler{},
		}

	*/
	router := map[string]*Tree{}
	router["GET"] = NewTree()
	router["POST"] = NewTree()
	router["PUT"] = NewTree()
	router["DELETE"] = NewTree()
	return &Core{router: router}
}

func (c *Core) Use(middlewares ...ControllerHandler) {
	c.middlewares = append(c.middlewares, middlewares...)
}

func (c *Core) Get(url string, handler ...ControllerHandler) {
	//c.router[url] = handler
	allHandlers := append(c.middlewares, handler...)
	if err := c.router["GET"].AddRouter(url, allHandlers); err != nil {
		log.Fatal("add router error: ", err)
	}
}

func (c *Core) Post(url string, handler ...ControllerHandler) {
	allHandlers := append(c.middlewares, handler...)
	if err := c.router["POST"].AddRouter(url, allHandlers); err != nil {
		log.Fatal("add router error: ", err)
	}
}

func (c *Core) Put(url string, handler ...ControllerHandler) {
	allHandlers := append(c.middlewares, handler...)
	if err := c.router["PUT"].AddRouter(url, allHandlers); err != nil {
		log.Fatal("add router error: ", err)
	}
}

func (c *Core) Delete(url string, handler ...ControllerHandler) {
	allHandlers := append(c.middlewares, handler...)
	if err := c.router["DELETE"].AddRouter(url, allHandlers); err != nil {
		log.Fatal("add router error:", err)
	}
}

func (c *Core) Group(prefix string) IGroup {
	return NewGroup(c, prefix)
}

//func (c *Core) FindRouteByRequest(request *http.Request) []ControllerHandler {
func (c *Core) FindRouteByRequest(request *http.Request) *node {
	uri := request.URL.Path
	method := request.Method
	upperMethod := strings.ToUpper(method)

	if methodHandlers, ok := c.router[upperMethod]; ok {
		//return methodHandlers.FindHandler(uri)
		return methodHandlers.root.matchNode(uri)
	}
	return nil
}

func (c *Core) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	log.Println("core.ServeHTTP")
	ctx := NewContext(r, w)

	//router := c.router["foo"]
	//handlers := c.FindRouteByRequest(r)
	node := c.FindRouteByRequest(r)
	//if handlers == nil {
	if node == nil {
		ctx.SetStatus(404).Json("not found")
		return
	}
	log.Println("core.Router")
	ctx.SetHandlers(node.handlers)

	//router(ctx)
	//ctx.SetHandlers(handlers)
	/*
		if err := router(ctx); err != nil {
			ctx.Json(500, "inner error")
		}

	*/
	params := node.parseParamsFromEndNode(r.URL.Path)
	ctx.SetParams(params)

	if err := ctx.Next(); err != nil {
		ctx.SetStatus(500).Json("inner error")
		return
	}
}
