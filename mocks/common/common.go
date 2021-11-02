package common

import (
	"fmt"
	"github.com/defernest/dmarket-go/dmarket"
	"io"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

type EndpointBehavior struct {
	httpMethod   string
	relativePath string
	handlerFunc  gin.HandlerFunc
}

func NewEndpointBehavior(httpMethod string, relativePath string, handlerFunc gin.HandlerFunc) *EndpointBehavior {
	return &EndpointBehavior{httpMethod: httpMethod, relativePath: relativePath, handlerFunc: handlerFunc}
}

func (e EndpointBehavior) Endpoint() (httpMethod string, relativePath string, handler gin.HandlerFunc) {
	return e.httpMethod, e.relativePath, e.handlerFunc
}

func MustReturnStatusOK(method, path string) *EndpointBehavior {
	return &EndpointBehavior{
		httpMethod:   method,
		relativePath: path,
		handlerFunc: func(context *gin.Context) {
			if context.Request.Body != http.NoBody {
				body, err := io.ReadAll(context.Request.Body)
				if err != nil {
					panic(fmt.Errorf("can`t read request body %w", err))
				}
				context.String(http.StatusOK, string(body))
			} else {
				context.String(http.StatusOK, "{}")
			}
		},
	}
}

func MustReturnBadBody(method, path string) *EndpointBehavior {
	return &EndpointBehavior{
		httpMethod:   method,
		relativePath: path,
		handlerFunc: func(context *gin.Context) {
			context.String(http.StatusOK, "{")
		},
	}
}

func MustReturnHTTPError(method, path string, errcode int) *EndpointBehavior {
	if errcode < 400 || errcode > 511 {
		panic(fmt.Errorf("the error code must be greater than 400 and less than 511"))
	}
	return &EndpointBehavior{
		httpMethod:   method,
		relativePath: path,
		handlerFunc: func(context *gin.Context) {
			context.String(dmarket.ErrorRepresentation{Code: errcode}.String())
		},
	}
}

func MustDoLongResponse(method, path string, sleep int) *EndpointBehavior {
	return &EndpointBehavior{
		httpMethod:   method,
		relativePath: path,
		handlerFunc: func(context *gin.Context) {
			time.Sleep(time.Duration(sleep) * time.Second)
			context.String(http.StatusOK, "{}")
		},
	}
}

func MustReturnDoBadRepresentationError(method, path string, code int) *EndpointBehavior {
	if code < 400 {
		panic("do bad representation error code must be more than 400")
	}
	return &EndpointBehavior{
		httpMethod:   method,
		relativePath: path,
		handlerFunc: func(context *gin.Context) {
			context.String(code, "{")
		},
	}
}
