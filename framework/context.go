package framework

import (
	"context"
	"net/http"
	"sync"
	"time"
)

type Context struct {
	request        *http.Request
	responseWriter http.ResponseWriter
	ctx            context.Context
	handler        ControllerHandler

	hasTimeout bool
	writerMux  *sync.Mutex

	handlers []ControllerHandler
	index    int

	params map[string]string
}

func NewContext(r *http.Request, w http.ResponseWriter) *Context {
	return &Context{
		request:        r,
		responseWriter: w,
		ctx:            r.Context(),
		writerMux:      &sync.Mutex{},
		index:          -1,
	}
}

// base func
func (c *Context) WriterMux() *sync.Mutex {
	return c.writerMux
}

func (c *Context) GetRequest() *http.Request {
	return c.request
}

func (c *Context) GetResponse() http.ResponseWriter {
	return c.responseWriter
}

func (c *Context) SetHasTimeout() {
	c.hasTimeout = true
}

func (c *Context) HasTimeout() bool {
	return c.hasTimeout
}

func (c *Context) SetHandlers(handlers []ControllerHandler) {
	c.handlers = handlers
}

func (c *Context) SetParams(params map[string]string) {
	c.params = params
}

func (c *Context) Next() error {
	c.index++
	if c.index < len(c.handlers) {
		if err := c.handlers[c.index](c); err != nil {
			return err
		}
	}
	return nil
}

// impl context
func (c *Context) BaseContext() context.Context {
	return c.request.Context()
}

func (c *Context) Deadline() (deadline time.Time, ok bool) {
	return c.BaseContext().Deadline()
}

func (c *Context) Done() <-chan struct{} {
	return c.BaseContext().Done()
}

func (c *Context) Err() error {
	return c.BaseContext().Err()
}

func (c *Context) Value(key interface{}) interface{} {
	return c.BaseContext().Value(key)
}

/*
// request function (URL Query)
func (c *Context) QueryAll() map[string][]string {
	if c.request != nil {
		return map[string][]string(c.request.URL.Query())
	}
	return map[string][]string{}
}

func (c *Context) QueryInt(key string, def int) int {
	params := c.QueryAll()
	if vals, ok := params[key]; ok {
		valsLen := len(vals)
		if valsLen > 0 {
			intval, err := strconv.Atoi(vals[valsLen-1])
			if err != nil {
				return def
			}
			return intval
		}
	}
	return def
}

func (c *Context) QueryString(key string, def string) string {
	params := c.QueryAll()
	if vals, ok := params[key]; ok {
		valsLen := len(vals)
		if valsLen > 0 {
			return vals[valsLen-1]
		}
		return def
	}
	return def
}

func (c *Context) QueryArray(key string, def []string) []string {
	params := c.QueryAll()
	if vals, ok := params[key]; ok {
		return vals
	}
	return def
}

// form post
func (c *Context) FormAll() map[string][]string {
	if c.request != nil {
		return map[string][]string(c.request.PostForm)
	}
	return map[string][]string{}
}

func (c *Context) FormInt(key string, def int) int {
	params := c.FormAll()
	if vals, ok := params[key]; ok {
		valsLen := len(vals)
		if valsLen > 0 {
			intval, err := strconv.Atoi(vals[valsLen-1])
			if err != nil {
				return intval
			}
			return def
		}
	}
	return def
}

func (c *Context) FormString(key string, def string) string {
	params := c.FormAll()
	if vals, ok := params[key]; ok {
		valsLen := len(vals)
		if valsLen > 0 {
			return vals[valsLen-1]
		}
		return def
	}
	return def
}

func (c *Context) FormArray(key string, def []string) []string {
	params := c.FormAll()
	if vals, ok := params[key]; ok {
		return vals
	}
	return def
}

// request application/json
func (c *Context) BindJson(obj interface{}) error {
	if c.request != nil {
		body, err := ioutil.ReadAll(c.request.Body)
		if err != nil {
			return err
		}
		c.request.Body = ioutil.NopCloser(bytes.NewBuffer(body))

		err = json.Unmarshal(body, &obj)
		if err != nil {
			return err
		}
	} else {
		return errors.New("ctx.request empty")
	}
	return nil
}

// response
func (c *Context) Json(status int, obj interface{}) error {
	if c.HasTimeout() {
		return nil
	}

	c.responseWriter.Header().Set("Content-Type", "application/json")
	c.responseWriter.WriteHeader(status)
	byt, err := json.Marshal(obj)
	if err != nil {
		c.responseWriter.WriteHeader(500)
		return err
	}
	c.responseWriter.Write(byt)
	return nil
}

func (c *Context) HTML(status int, obj interface{}, template string) error {
	return nil
}

func (c *Context) Text(status int, obj string) error {
	return nil
}

*/
