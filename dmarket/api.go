package dmarket

import (
	"context"
	"crypto/ed25519"
	"encoding/hex"
	"fmt"
	"github.com/hashicorp/go-multierror"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"time"

	"golang.org/x/time/rate"
)

type defaultClient struct {
	http                  *http.Client
	rateLimit             *rate.Limiter
	baseURL               *url.URL
	publicKey, privateKey string
}

func (c defaultClient) Get(endpoint string) (Response, error) {
	req, err := http.NewRequest(http.MethodGet, endpoint, http.NoBody)
	if err != nil {
		return Response{}, err
	}
	return c.Do(req)
}

func (c defaultClient) Post(endpoint string, body io.Reader) (Response, error) {
	req, err := http.NewRequest(http.MethodPost, endpoint, body)
	if err != nil {
		return Response{}, err
	}
	defer func() {
		err := req.Body.Close()
		if err != nil {
			panic("failed when defer resp body close")
		}
	}()
	return c.Do(req)
}

func (c defaultClient) Delete(endpoint string, body io.Reader) (Response, error) {
	req, err := http.NewRequest(http.MethodDelete, endpoint, body)
	if err != nil {
		return Response{}, err
	}
	defer func() {
		err := req.Body.Close()
		if err != nil {
			panic("failed when defer resp body close")
		}
	}()
	return c.Do(req)
}

func (c defaultClient) Patch(endpoint string, body io.Reader) (Response, error) {
	req, err := http.NewRequest(http.MethodPatch, endpoint, body)
	if err != nil {
		return Response{}, err
	}
	defer func() {
		err := req.Body.Close()
		if err != nil {
			panic("failed when defer resp body close")
		}
	}()
	return c.Do(req)
}

/*
sign a request to the Dmarket Items API with a private key, setting X-Api-Key, X-Sign-Date X-Request-Sign headers for this request

To make a signature (X-Request-Sign), take the following steps:
	1. Build non-signed string formula:
	   (HTTP Method) + (Route path + HTTP query params) + (body string) + (timestamp)
    2. After you’ve created a non-signed string with a default concatenation method,
       sign it with ed25519 using you secret key.
    3. Encode the result string with hex
*/
func (c defaultClient) sign(req *http.Request) error {
	timestamp := strconv.FormatInt(time.Now().UTC().Unix(), 10)
	b, err := hex.DecodeString(c.privateKey)
	if err != nil {
		return fmt.Errorf("api: decode private key error: %w", err)
	}
	var privateKey [64]byte
	copy(privateKey[:], b[:64])
	req.Header.Set("X-Sign-Date", timestamp)
	req.Header.Set("X-Request-Sign", "dmar ed25519 "+hex.EncodeToString(ed25519.Sign(privateKey[:], []byte(req.Method+req.URL.RequestURI()+timestamp))))
	req.Header.Set("X-Api-Key", c.publicKey)
	return nil
}

/*
Do performs a request to the Dmarket Items API
*/
func (c *defaultClient) Do(req *http.Request) (response Response, errs error) {
	defer func() {
		if err := recover(); err != nil {
			errs = multierror.Append(errs, fmt.Errorf("unexpected error when Do request - abort!\n\terror: %s", err))
		}
	}()
	err := c.rateLimit.Wait(context.TODO())
	if err != nil {
		return Response{}, fmt.Errorf("api: request rate limiter error: %w", err)
	}
	req.URL = c.baseURL.ResolveReference(req.URL)
	err = c.sign(req)
	if err != nil {
		return Response{}, fmt.Errorf("api: new request sign error: %w", err)
	}
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Content-Type", "application/json")
	resp, err := c.http.Do(req)
	if err != nil {
		return Response{}, fmt.Errorf("api: client Do request error: %w", err)
	}
	defer func() {
		if resp.Body != nil {
			if closeErr := resp.Body.Close(); closeErr != nil {
				errs = multierror.Append(errs, fmt.Errorf("failed when defer resp body close: %w", closeErr))
			}
		}
	}()
	response.Status = resp.Status
	response.StatusCode = resp.StatusCode
	response.ContentLength = resp.ContentLength
	response.Request = resp.Request
	_, err = response.ReadFrom(resp.Body)
	if err != nil {
		return Response{}, fmt.Errorf("api: read responce body error: %w", err)
	}
	return response, nil
}
