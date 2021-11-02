package dmarket

import (
	"fmt"
	"net/http"
	"net/url"
	"time"

	"golang.org/x/time/rate"
)

type Client struct {
	DefaultClient *defaultClient

	Exchange *Exchange
}

type errorBadKeys struct {
	public, private string
}

func (e errorBadKeys) Error() string {
	return fmt.Sprintf("wrong public or private key len\n"+
		"public key: %1s len: %d must be 64\n"+
		"private key: %s len: %d must be 128\n", e.public, len(e.public), e.private, len(e.private))
}

// NewClient create a new Dmarket API client
func NewClient(baseURL, publicKey, privateKey string) (*Client, error) {
	if len(publicKey) != 64 || len(privateKey) != 128 {
		return nil, errorBadKeys{public: publicKey, private: privateKey}
	}
	base, err := url.Parse(baseURL)
	if err != nil || base.Hostname() == "" {
		return nil, fmt.Errorf("baseURL url parsing error: walid hostname format [scheme:][//[userinfo@]baseURL]")
	}
	c := &Client{
		DefaultClient: &defaultClient{
			http:       &http.Client{Timeout: 10 * time.Second},
			rateLimit:  rate.NewLimiter(rate.Every(125*time.Millisecond), 8),
			baseURL:    base,
			publicKey:  publicKey,
			privateKey: privateKey,
		},
	}
	c.Exchange = NewExchange(c.DefaultClient)
	return c, nil
}
