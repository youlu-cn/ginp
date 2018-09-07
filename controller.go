package ginp

import (
	"net/http"
	"reflect"
	"strings"

	"github.com/gin-gonic/gin"
)

type Controller interface {
	Group() string
}

type Response struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

func ErrorResponse(code int, err error) *Response {
	msg := ""
	if err != nil {
		msg = err.Error()
	}
	return &Response{
		Code:    code,
		Message: msg,
	}
}

func DataResponse(data interface{}) *Response {
	return &Response{
		Data: data,
	}
}

type Request struct {
	*http.Request
	ctx *gin.Context
}

func (req *Request) Bind(obj interface{}) error {
	return req.ctx.Bind(obj)
}

func (req *Request) BindJSON(obj interface{}) error {
	return req.ctx.BindJSON(obj)
}

func handleRequest(obj interface{}) gin.HandlerFunc {
	return func(c *gin.Context) {
		req := &Request{
			Request: c.Request,
			ctx:     c,
		}

		// parse path
		path := strings.Trim(c.Request.URL.Path, "/")
		urlParts := strings.Split(path, "/")
		parts := make([]string, 0, len(urlParts)-1)
		for _, v := range urlParts[1:] {
			parts = append(parts, strings.Title(v))
		}
		// get method
		controller := reflect.ValueOf(obj)
		handler := strings.Join(parts, "") + "Handler"
		action := strings.Title(strings.ToLower(c.Request.Method)) + "Action"
		method := controller.MethodByName(handler + action)
		if !method.IsValid() {
			method = controller.MethodByName(handler)
			if !method.IsValid() {
				goto NotFound
			}
		}
		// check method
		if method.Type().NumIn() != 1 || method.Type().NumOut() != 1 {
			goto NotFound
		}
		// call
		if rets := method.Call([]reflect.Value{reflect.ValueOf(req)}); len(rets) != 1 {
			goto NotFound
		} else if resp, ok := rets[0].Interface().(*Response); !ok {
			goto NotFound
		} else {
			c.JSON(http.StatusOK, resp)
			return
		}

	NotFound:
		c.AbortWithStatus(http.StatusNotFound)
	}
}
