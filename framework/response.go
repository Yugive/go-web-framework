package framework

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
	"html/template"
	"net/http"
	"net/url"
)

type IResponse interface {
	Json(obj interface{}) IResponse

	Jsonp(obj interface{}) IResponse

	Xml(obj interface{}) IResponse

	Html(template string, obj interface{}) IResponse

	Text(format string, values ...interface{}) IResponse

	Redirect(path string) IResponse

	SetHeader(key string, val string) IResponse

	SetCookie(key string, val string, maxAge int, path, domain string, secure, httpOnly bool) IResponse

	SetStatus(code int) IResponse
}

func (c *Context) SetHeader(key string, val string) IResponse {
	c.responseWriter.Header().Add(key, val)
	return c
}

func (c *Context) SetStatus(code int) IResponse {
	c.responseWriter.WriteHeader(code)
	return c
}

func (c *Context) SetOkStatus() IResponse {
	c.responseWriter.WriteHeader(http.StatusOK)
	return c
}

func (c *Context) Redirect(path string) IResponse {
	http.Redirect(c.responseWriter, c.request, path, http.StatusMovedPermanently)
	return c
}

func (c *Context) Text(format string, values ...interface{}) IResponse {
	out := fmt.Sprintf(format, values...)
	c.SetHeader("Content-type", "application/text")
	c.responseWriter.Write([]byte(out))
	return c
}

func (c *Context) SetCookie(key, val string, maxAge int, path, domain string, secure bool, httpOnly bool) IResponse {
	if path == "" {
		path = "/"
	}
	http.SetCookie(c.responseWriter, &http.Cookie{
		Name:     key,
		Value:    url.QueryEscape(val),
		MaxAge:   maxAge,
		Path:     path,
		Domain:   domain,
		SameSite: 1,
		Secure:   secure,
		HttpOnly: httpOnly,
	})
	return c
}

func (c *Context) Json(obj interface{}) IResponse {
	byt, err := json.Marshal(obj)
	if err != nil {
		return c.SetStatus(http.StatusInternalServerError)
	}
	c.SetHeader("Content-Type", "application/json")
	c.responseWriter.Write(byt)
	return c
}

func (c *Context) Jsonp(obj interface{}) IResponse {
	callbackFunc, _ := c.QueryString("callback", "callback_function")
	c.SetHeader("Content-Type", "application/javascript")
	callback := template.JSEscapeString(callbackFunc)

	_, err := c.responseWriter.Write([]byte(callback))
	if err != nil {
		return c
	}

	_, err = c.responseWriter.Write([]byte("("))
	if err != nil {
		return c
	}

	ret, err := json.Marshal(obj)
	if err != nil {
		return c
	}

	_, err = c.responseWriter.Write(ret)
	if err != nil {
		return c
	}

	_, err = c.responseWriter.Write([]byte(")"))
	if err != nil {
		return c
	}

	return c
}

func (c *Context) Xml(obj interface{}) IResponse {
	byt, err := xml.Marshal(obj)
	if err != nil {
		return c.SetStatus(http.StatusInternalServerError)
	}
	c.SetHeader("Content-Type", "application/html")
	c.responseWriter.Write(byt)
	return c
}

func (c *Context) Html(file string, obj interface{}) IResponse {
	t, err := template.New("output").ParseFiles(file)
	if err != nil {
		return c
	}

	if err := t.Execute(c.responseWriter, obj); err != nil {
		return c
	}
	c.SetHeader("Content-Type", "text/html")
	return c
}
