package mocks

import (
	"context"
	"crypto/ed25519"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"github.com/defernest/dmarket-go/dmarket"
	"io"
	"net/http"
	"strconv"
	"time"

	"golang.org/x/time/rate"
)

type dmarketClient struct {
	server    *DmarketServer
	rateLimit *rate.Limiter
}

func (c dmarketClient) Get(endpoint string) (dmarket.Response, error) {
	req, err := http.NewRequest(http.MethodGet, c.server.URL()+endpoint, http.NoBody)
	if err != nil {
		return dmarket.Response{}, err
	}
	return c.Do(req)
}

func (c dmarketClient) Post(endpoint string, body io.Reader) (dmarket.Response, error) {
	req, err := http.NewRequest(http.MethodGet, c.server.URL()+endpoint, body)
	if err != nil {
		return dmarket.Response{}, err
	}
	return c.Do(req)
}

func (c dmarketClient) Delete(endpoint string, body io.Reader) (dmarket.Response, error) {
	panic(fmt.Errorf("implement me %s, %v", endpoint, body))
}

func (c dmarketClient) Patch(endpoint string, body io.Reader) (dmarket.Response, error) {
	panic(fmt.Errorf("implement me %s, %v", endpoint, body))
}

func (c dmarketClient) Do(req *http.Request) (dmarket.Response, error) {
	err := c.rateLimit.Wait(context.TODO())
	if err != nil {
		return dmarket.Response{}, fmt.Errorf("api: request rate limiter error: %w", err)
	}
	timestamp := strconv.Itoa(int(time.Now().UTC().Unix()))
	signature, err := c.server.sign(req.Method, req.URL.RequestURI(), timestamp)
	if err != nil {
		return dmarket.Response{}, fmt.Errorf("api mock: new request sign error: %w", err)
	}
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Sign-Date", timestamp)
	req.Header.Set("X-Request-Sign", "dmar ed25519 "+signature)
	req.Header.Set("X-Api-Key", c.server.PublicKey)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return dmarket.Response{}, fmt.Errorf("api mock: do request sign error: %w", err)
	}
	r := dmarket.Response{
		Status:        resp.Status,
		StatusCode:    resp.StatusCode,
		ContentLength: resp.ContentLength,
		Request:       resp.Request,
	}
	_, err = r.ReadFrom(resp.Body)
	if err != nil {
		return dmarket.Response{}, fmt.Errorf("api mock: read responce body error: %w", err)
	}
	return r, nil
}

func (s *DmarketServer) generateKeys() {
	pub, private, err := ed25519.GenerateKey(rand.Reader)
	if err != nil {
		panic(err)
	}
	s.PublicKey = hex.EncodeToString(pub)
	s.PrivareKey = hex.EncodeToString(private)
}

/*
sign creates a signature for the X-Request-Sign header

To make a signature, take the following steps:
	1. Build non-signed string formula:
	   (HTTP Method) + (Route path + HTTP query params) + (body string) + (timestamp) )
    2. After you’ve created a non-signed string with a default concatenation method,
       sign it with ed25519 using you secret key.
    3. Encode the result string with hex
*/
func (s DmarketServer) sign(method, path, timestamp string) (string, error) {
	b, err := hex.DecodeString(s.PrivareKey)
	if err != nil {
		return "", fmt.Errorf("api: sign decode string error: %w", err)
	}
	var privateKey [64]byte
	copy(privateKey[:], b[:64])
	sign := hex.EncodeToString(ed25519.Sign(privateKey[:], []byte(method+path+timestamp)))
	return sign, nil
}

func (s DmarketServer) wrongGet(wrongTimestamp, wrongPath, wrongPublic, wrongSign bool) (*http.Response, error) {
	pub := s.PublicKey
	path := s.URL()
	timestamp := strconv.Itoa(int(time.Now().UTC().Unix()))
	if wrongTimestamp {
		timestamp = strconv.Itoa(int(time.Now().Add(1 * time.Hour).UTC().Unix()))
	}
	if wrongPath {
		path = s.URL() + "/err"
	}
	signature, err := s.sign(http.MethodGet, s.URL(), timestamp)
	if err != nil {
		return nil, fmt.Errorf("api: new request sign error: %w", err)
	}
	req, err := http.NewRequest(http.MethodGet, path, http.NoBody)
	if err != nil {
		return nil, fmt.Errorf("error creating dmarket request (%w)", err)
	}
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Sign-Date", timestamp)
	if wrongSign {
		signature = signature[10:]
	}
	req.Header.Set("X-Request-Sign", "dmar ed25519 "+signature)
	if wrongPublic {
		pub = pub[2:]
	}
	req.Header.Set("X-Api-Key", s.PublicKey)
	return http.DefaultClient.Do(req)
}
