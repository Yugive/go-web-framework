package framework

import (
	"bytes"
	"encoding/json"
	"encoding/xml"
	"errors"
	"github.com/spf13/cast"
	"io/ioutil"
	"mime/multipart"
)

const defaultMultipartMemory = 32 << 20

type IRequest interface {
	QueryInt(key string, def int) (int, bool)
	QueryInt64(key string, def int64) (int64, bool)
	QueryFloat64(key string, def float64) (float64, bool)
	QueryFloat32(key string, def float32) (float32, bool)
	QueryBool(key string, def bool) (bool, bool)
	QueryString(key string, def string) (string, bool)
	QueryStringSlice(key string, def []string) ([]string, bool)
	Query(key string) interface{}

	ParamInt(key string, def int) (int, bool)
	ParamInt64(key string, def int64) (int64, bool)
	ParamFloat64(key string, def float64) (float64, bool)
	ParamFloat32(key string, def float32) (float32, bool)
	ParamBool(key string, def bool) (bool, bool)
	ParamString(key string, def string) (string, bool)
	Param(key string) interface{}

	FormInt(key string, def int) (int, bool)
	FormInt64(key string, def int64) (int64, bool)
	FormFloat64(key string, def float64) (float64, bool)
	FormFloat32(key string, def float32) (float32, bool)
	FormBool(key string, def bool) (bool, bool)
	FormString(key string, def string) (string, bool)
	FormStringSlice(key string, def []string) ([]string, bool)
	FormFile(key string) (*multipart.FileHeader, error)
	Form(key string) interface{}

	BindJson(obj interface{}) error

	BindXml(obj interface{}) error

	GetRawData() ([]byte, error)

	Uri() string
	Method() string
	Host() string
	ClientIp() string

	Headers() map[string][]string
	Header(key string) (string, bool)

	Cookies() map[string]string
	Cookie(key string) (string, bool)
}

func (c *Context) QueryAll() map[string][]string {
	if c.request != nil {
		return map[string][]string(c.request.URL.Query())
	}
	return map[string][]string{}
}

func (c *Context) QueryInt(key string, def int) (int, bool) {
	params := c.QueryAll()
	if vals, ok := params[key]; ok {
		if len(vals) > 0 {
			return cast.ToInt(vals[0]), true
		}
	}
	return def, false
}

func (c *Context) QueryInt64(key string, def int64) (int64, bool) {
	params := c.QueryAll()
	if vals, ok := params[key]; ok {
		if len(vals) > 0 {
			return cast.ToInt64(vals[0]), true
		}
	}
	return def, false
}

func (c *Context) QueryFloat64(key string, def float64) (float64, bool) {
	params := c.QueryAll()
	if vals, ok := params[key]; ok {
		if len(vals) > 0 {
			return cast.ToFloat64(vals[0]), true
		}
	}
	return def, false
}

func (c *Context) QueryFloat32(key string, def float32) (float32, bool) {
	params := c.QueryAll()
	if vals, ok := params[key]; ok {
		if len(vals) > 0 {
			return cast.ToFloat32(vals[0]), true
		}
	}
	return def, false
}

func (c *Context) QueryString(key string, def string) (string, bool) {
	params := c.QueryAll()
	if vals, ok := params[key]; ok {
		if len(vals) > 0 {
			return vals[0], true
		}
	}
	return def, false
}

func (c *Context) QueryStringSlice(key string, def []string) ([]string, bool) {
	params := c.QueryAll()
	if vals, ok := params[key]; ok {
		return vals, true
	}
	return def, false
}

func (c *Context) QueryBool(key string, def bool) (bool, bool) {
	params := c.QueryAll()
	if vals, ok := params[key]; ok {
		if len(vals) > 0 {
			return cast.ToBool(vals[0]), true
		}
	}
	return def, false
}

func (c *Context) Query(key string) interface{} {
	params := c.QueryAll()
	if vals, ok := params[key]; ok {
		return vals[0]
	}
	return nil
}

func (c *Context) Param(key string) interface{} {
	if c.params != nil {
		if val, ok := c.params[key]; ok {
			return val
		}
	}
	return nil
}

func (c *Context) ParamInt(key string, def int) (int, bool) {
	if val := c.Param(key); val != nil {
		return cast.ToInt(val), true
	}
	return def, false
}

func (c *Context) ParamInt64(key string, def int64) (int64, bool) {
	if val := c.Param(key); val != nil {
		return cast.ToInt64(val), true
	}
	return def, false
}

func (c *Context) ParamFloat64(key string, def float64) (float64, bool) {
	if val := c.Param(key); val != nil {
		return cast.ToFloat64(val), true
	}
	return def, false
}

func (c *Context) ParamBool(key string, def bool) (bool, bool) {
	if val := c.Param(key); val != nil {
		return cast.ToBool(val), true
	}
	return def, false
}

func (c *Context) ParamString(key string, def string) (string, bool) {
	if val := c.Param(key); val != nil {
		return cast.ToString(val), true
	}
	return def, false
}

func (c *Context) FormAll() map[string][]string {
	if c.request != nil {
		c.request.ParseForm()
		return map[string][]string(c.request.PostForm)
	}
	return map[string][]string{}
}

func (c *Context) FormInt(key string, def int) (int, bool) {
	params := c.FormAll()
	if vals, ok := params[key]; ok {
		if len(vals) > 0 {
			return cast.ToInt(vals[0]), true
		}
	}
	return def, false
}

func (c *Context) FormInt64(key string, def int64) (int64, bool) {
	params := c.FormAll()
	if vals, ok := params[key]; ok {
		if len(vals) > 0 {
			return cast.ToInt64(vals[0]), true
		}
	}
	return def, false
}

func (c *Context) FormFloat64(key string, def float64) (float64, bool) {
	params := c.FormAll()
	if vals, ok := params[key]; ok {
		if len(vals) > 0 {
			return cast.ToFloat64(vals[0]), true
		}
	}

	return def, false
}

func (c *Context) FormFloat32(key string, def float32) (float32, bool) {
	params := c.FormAll()
	if vals, ok := params[key]; ok {
		if len(vals) > 0 {
			return cast.ToFloat32(vals[0]), true
		}
	}
	return def, false
}

func (c *Context) FormBool(key string, def bool) (bool, bool) {
	params := c.FormAll()
	if vals, ok := params[key]; ok {
		if len(vals) > 0 {
			return cast.ToBool(vals[0]), true
		}
	}
	return def, false
}

func (c *Context) FormString(key string, def string) (string, bool) {
	params := c.FormAll()
	if vals, ok := params[key]; ok {
		if len(vals) > 0 {
			return vals[0], true
		}
	}
	return def, false
}

func (c *Context) FormStringSlice(key string, def []string) ([]string, bool) {
	params := c.FormAll()
	if vals, ok := params[key]; ok {
		return vals, true
	}
	return def, false
}

func (c *Context) FormFile(key string) (*multipart.FileHeader, error) {
	if c.request.MultipartForm == nil {
		if err := c.request.ParseMultipartForm(defaultMultipartMemory); err != nil {
			return nil, err
		}
	}
	f, fh, err := c.request.FormFile(key)
	if err != nil {
		return nil, err
	}
	f.Close()
	return fh, err
}

func (c *Context) Form(key string) interface{} {
	params := c.FormAll()
	if vals, ok := params[key]; ok {
		if len(vals) > 0 {
			return vals[0]
		}
	}
	return nil
}

func (c *Context) BindJson(obj interface{}) error {
	if c.request != nil {
		body, err := ioutil.ReadAll(c.request.Body)
		if err != nil {
			return err
		}
		c.request.Body = ioutil.NopCloser(bytes.NewBuffer(body))
		err = json.Unmarshal(body, obj)
		if err != nil {
			return err
		}
	} else {
		return errors.New("ctx request empty")
	}
	return nil
}

func (c *Context) BindXml(obj interface{}) error {
	if c.request != nil {
		body, err := ioutil.ReadAll(c.request.Body)
		if err != nil {
			return err
		}

		c.request.Body = ioutil.NopCloser(bytes.NewBuffer(body))

		err = xml.Unmarshal(body, obj)
		if err != nil {
			return err
		}
	} else {
		return errors.New("ctx request empty")
	}
	return nil
}

func (c *Context) GetRawData() ([]byte, error) {
	if c.request != nil {
		body, err := ioutil.ReadAll(c.request.Body)
		if err != nil {
			return nil, err
		}

		c.request.Body = ioutil.NopCloser(bytes.NewBuffer(body))
		return body, nil
	}
	return nil, errors.New("ctx request empty")
}

func (c *Context) Uri() string {
	return c.request.RequestURI
}

func (c *Context) Method() string {
	return c.request.Method
}

func (c *Context) Host() string {
	return c.request.URL.Host
}

func (c *Context) ClientIp() string {
	r := c.request
	ipAddress := r.Header.Get("X-Real-Ip")
	if ipAddress == "" {
		ipAddress = r.Header.Get("X-Forwarded-For")
	}
	if ipAddress == "" {
		ipAddress = r.RemoteAddr
	}
	return ipAddress
}

func (c *Context) Headers() map[string][]string {
	return map[string][]string(c.request.Header)
}

func (c *Context) Header(key string) (string, bool) {
	vals := c.request.Header.Values(key)
	if vals == nil || len(vals) < 0 {
		return "", false
	}
	return vals[0], true
}

func (c *Context) Cookies() map[string]string {
	cookies := c.request.Cookies()
	ret := map[string]string{}
	for _, cookie := range cookies {
		ret[cookie.Name] = cookie.Value
	}
	return ret
}

func (c *Context) Cookie(key string) (string, bool) {
	cookies := c.Cookies()
	if val, ok := cookies[key]; ok {
		return val, true
	}
	return "", false
}
