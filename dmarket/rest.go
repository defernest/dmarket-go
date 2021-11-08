package dmarket

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
)

type Requester interface {
	Get(endpoint string) (Response, error)
	Post(endpoint string, body io.Reader) (Response, error)
	Delete(endpoint string, body io.Reader) (Response, error)
	Patch(endpoint string, body io.Reader) (Response, error)
}

type ErrorRepresentation struct {
	Response Response
}

func (e ErrorRepresentation) Error() string {
	return fmt.Sprintf("dmarket API representation error: code %d: %s", e.Response.StatusCode, http.StatusText(e.Response.StatusCode))
}

func (e ErrorRepresentation) String() (int, string) {
	return e.Response.StatusCode, fmt.Sprintf("%d: %s", e.Response.StatusCode, http.StatusText(e.Response.StatusCode))
}

type RuntimeError struct {
	Err     string `json:"error"`
	Code    int    `json:"code"`
	Message string `json:"message"`
	Details []struct {
		TypeURL string `json:"type_url"`
		Value   string `json:"value"`
	} `json:"details"`
}

func (e RuntimeError) Error() string {
	var dv string
	for i, d := range e.Details {
		dv = dv + fmt.Sprintf("[%d](typeURL: %s, value: %s) ", i, d.TypeURL, d.Value)
	}
	return fmt.Sprintf("dmarket API runtime error::: err: %s [code: %d], message: %s details: %s",
		e.Err, e.Code, e.Message, dv)
}

type Response struct {
	Status        string
	StatusCode    int
	ContentLength int64
	Body          *bytes.Buffer
	Request       *http.Request
}

func (r *Response) ReadFrom(reader io.Reader) (n int64, err error) {
	r.Body = new(bytes.Buffer)
	return r.Body.ReadFrom(reader)
}
