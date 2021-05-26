package client

import (
	"context"
	"strings"
	"time"

	"github.com/valyala/fasthttp"
)

type HTTPMethod string

const (
	POST   HTTPMethod = "POST"
	GET    HTTPMethod = "GET"
	PUT    HTTPMethod = "PUT"
	DELETE HTTPMethod = "DELETE"
)

type HTTPRequest struct {
	URL         string
	Method      HTTPMethod
	ContentType string
	Headers     map[string]string
	Body        string
}

type HTTPResponse struct {
	Body       string
	Headers    map[string]string
	StatusCode int32
	Time       int64
}

type HTTPClient struct {
	Client      *fasthttp.Client
	baseURI     string
	contentType string
	timeout     time.Duration
}

type NewOpts struct {
	BaseURI            string
	DefaultContentType string
	Timeout            time.Duration
	Attemps            int
	TLSCert            string
}

type SendRequestOpts struct {
	StartTime *time.Time
	Request   *HTTPRequest
}

func NewHTTPClient(opts ...*NewOpts) *HTTPClient {
	opt := &NewOpts{
		DefaultContentType: "application/json",
		Timeout:            30 * time.Second,
		Attemps:            5,
	}

	if len(opts) > 0 {
		opt = opts[0]
	}

	c := &HTTPClient{
		Client:      &fasthttp.Client{},
		baseURI:     opt.BaseURI,
		contentType: opt.DefaultContentType,
		timeout:     opt.Timeout,
	}

	c.Client.MaxIdemponentCallAttempts = opt.Attemps

	return c
}

func (h *HTTPClient) SendRequest(ctx context.Context, opts *SendRequestOpts) (*HTTPResponse, error) {
	if opts.StartTime == nil {
		n := time.Now()
		opts.StartTime = &n
	}

	req := fasthttp.AcquireRequest()
	res := fasthttp.AcquireResponse()
	defer fasthttp.ReleaseRequest(req)
	defer fasthttp.ReleaseResponse(res)

	req.SetRequestURI(h.getURI(opts.Request.URL))
	req.Header.SetMethod(string(opts.Request.Method))

	cType := h.contentType
	if opts.Request.ContentType != "" {
		cType = opts.Request.ContentType
	}

	req.Header.SetContentType(cType)
	req.SetBodyString(opts.Request.Body)

	for key, value := range opts.Request.Headers {
		req.Header.Set(key, value)
	}

	if err := h.Client.DoTimeout(req, res, h.timeout); err != nil {
		return nil, err
	}

	endNow := time.Now()

	return &HTTPResponse{
		StatusCode: int32(res.StatusCode()),
		Body:       string(res.Body()),
		Headers:    mergeResponseHeaders(&res.Header),
		Time:       getResponseTime(opts.StartTime, &endNow),
	}, nil
}

func mergeResponseHeaders(h *fasthttp.ResponseHeader) map[string]string {
	headers := map[string]string{}

	h.VisitAll(func(key, value []byte) {
		headers[string(key)] = string(value)
	})

	return headers
}

func getResponseTime(start, end *time.Time) int64 {
	return end.Sub(*start).Milliseconds()
}

func (h *HTTPClient) getURI(rURI string) string {
	if strings.HasPrefix(rURI, "http") || strings.HasPrefix(rURI, "https") {
		return rURI
	}

	return h.baseURI + rURI
}
