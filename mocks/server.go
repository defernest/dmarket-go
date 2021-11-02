package mocks

import (
	"bytes"
	"crypto/ed25519"
	"encoding/hex"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/defernest/dmarket-go/dmarket"

	"golang.org/x/time/rate"

	"github.com/gin-gonic/gin"
)

type DmarketEndpoint interface {
	Endpoint() (httpMethod string, relativePath string, handler gin.HandlerFunc)
}

type DmarketServer struct {
	ts         *httptest.Server
	Client     *dmarketClient
	logs       *bytes.Buffer
	PrivareKey string
	PublicKey  string
}

func NewDmarketServer(endpoint DmarketEndpoint) DmarketServer {
	gin.SetMode(gin.TestMode)
	var logs bytes.Buffer
	gin.DefaultWriter = &logs

	router := gin.New()
	router.RedirectTrailingSlash = false
	router.Use(gin.LoggerWithFormatter(logger()), rateLimit(), checkHeaders(), dmarketAuth())
	router.NoRoute(noRoute())
	router.Handle(endpoint.Endpoint())

	s := DmarketServer{ts: httptest.NewServer(router), logs: &logs}
	s.generateKeys()
	s.Client = &dmarketClient{&s, rate.NewLimiter(5, 5)}
	return s
}

func (s DmarketServer) URL() string {
	return s.ts.URL
}

func (s DmarketServer) Logs() {
	_, err := fmt.Fprintln(os.Stderr, s.logs)
	if err != nil {
		fmt.Println("error writing dmarket server logs to stderr")
	}
}

func (s DmarketServer) Close() {
	defer s.ts.Close()
}

func rateLimit() gin.HandlerFunc {
	limiter := rate.NewLimiter(10, 5)
	return func(context *gin.Context) {
		if !limiter.Allow() {
			context.String(dmarket.ErrorRepresentation{Code: http.StatusTooManyRequests}.String())
			context.Abort()
		}
	}
}

func logger() gin.LogFormatter {
	return func(params gin.LogFormatterParams) string {
		return fmt.Sprintf("%s[GIN %s]%s -> %s %s[%d]%s\n%s",
			params.MethodColor(),
			params.Method,
			params.ResetColor(),
			params.Request.RequestURI,
			params.StatusCodeColor(),
			params.StatusCode,
			params.ResetColor(),
			params.ErrorMessage)
	}
}

func noRoute() gin.HandlerFunc {
	return func(context *gin.Context) {
		context.String(dmarket.ErrorRepresentation{Code: http.StatusNotFound}.String())
		context.AbortWithError(http.StatusNotFound, fmt.Errorf("no route to path '%s'", context.Request.RequestURI))
	}
}

func checkHeaders() gin.HandlerFunc {
	return func(context *gin.Context) {
		pub, err := hex.DecodeString(context.GetHeader("X-Api-Key"))
		if err != nil || len(pub) != 32 {
			context.String(dmarket.ErrorRepresentation{Code: http.StatusBadRequest}.String())
			context.AbortWithError(http.StatusBadRequest, errors.New("X-Api-Key error: decode error or len not equal 32"))
			return
		}
		context.Set("X-Api-Key", pub)

		signHeader := context.GetHeader("X-Request-Sign")
		containDMAR := strings.HasPrefix(signHeader, "dmar ed25519 ")
		sign, err := hex.DecodeString(strings.TrimPrefix(signHeader, "dmar ed25519 "))
		if err != nil || !containDMAR || len(sign) < 1 {
			context.String(dmarket.ErrorRepresentation{Code: http.StatusBadRequest}.String())
			context.AbortWithError(http.StatusBadRequest, errors.New("X-Request-Sign error: decode error, not contain dmar prefix or len < 1"))
			return
		}
		context.Set("X-Request-Sign", sign)

		signdate, err := strconv.Atoi(context.GetHeader("X-Sign-Date"))
		if err != nil || signdate <= 0 {
			context.String(dmarket.ErrorRepresentation{Code: http.StatusBadRequest}.String())
			context.AbortWithError(http.StatusBadRequest, errors.New("X-Sign-Date error: atoi error or len <= 0"))
			return
		}
		context.Set("X-Sign-Date", signdate)

		if context.GetHeader("Accept") != "application/json" || context.GetHeader("Content-Type") != "application/json" {
			context.String(dmarket.ErrorRepresentation{Code: http.StatusBadRequest}.String())
			context.AbortWithError(http.StatusBadRequest, errors.New("accept or content-Type headers not equal 'application/json'"))
		}
		return
	}
}

func dmarketAuth() gin.HandlerFunc {
	return func(context *gin.Context) {
		timestamp := time.Now().Add(2 * time.Minute).UTC().Unix()
		signdate := int64(context.GetInt("X-Sign-Date"))
		pub, _ := context.Get("X-Api-Key")
		sign, _ := context.Get("X-Request-Sign")
		msg := []byte(context.Request.Method + context.Request.URL.String() + strconv.FormatInt(signdate, 10))
		if signdate > timestamp || !ed25519.Verify(pub.([]byte), msg, sign.([]byte)) {
			_ = context.Error(errors.New("auth error")).
				SetMeta(
					fmt.Sprintf("X-Sign-Date correct: [%t] (false OK) | ed25519 verification: [%t] (true OK)",
						signdate > timestamp, ed25519.Verify(pub.([]byte), msg, sign.([]byte))),
				)
			context.String(dmarket.ErrorRepresentation{Code: http.StatusUnauthorized}.String())
			context.AbortWithStatus(http.StatusUnauthorized)
		}
		return
	}
}
