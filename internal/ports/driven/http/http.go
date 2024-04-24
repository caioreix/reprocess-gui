package http

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"reprocess-gui/internal/ports/driven/cache"
	"reprocess-gui/internal/utils"
)

const defaultCallerSkip = 4

// Service represents an HTTP client service.
type Service interface {
	Get(ctx context.Context, url string, payload io.Reader, headers map[string]string) (*Response, error)
	Post(ctx context.Context, url string, payload io.Reader, headers map[string]string) (*Response, error)

	Cached(ttl time.Duration, key ...string) Service
	SetCache(key ...string) Service
	Skip(size int) Service
}

// Response represents an HTTP response.
type Response struct {
	StatusCode int
	Body       []byte
}

type httpCache struct {
	client     cache.Service
	enabled    bool
	ttl        time.Duration
	key        string
	callerSkip int
}

type client struct {
	client *http.Client
	cache  *httpCache
}

// NewClient creates a new HTTP client with the specified timeout, transport, and cache.
func NewClient(timeout time.Duration, transport http.RoundTripper, cache cache.Service) Service {
	return &client{
		client: &http.Client{
			Timeout:   timeout,
			Transport: transport,
		},
		cache: &httpCache{
			client:     cache,
			callerSkip: defaultCallerSkip,
		},
	}
}

// Get performs an HTTP GET request.
func (c *client) Get(ctx context.Context, url string, payload io.Reader, headers map[string]string) (*Response, error) {
	res, err := c.do(ctx, "GET", url, payload, headers)
	if err != nil {
		return nil, err
	}

	return res, nil
}

// Post performs an HTTP POST request.
func (c *client) Post(ctx context.Context, url string, payload io.Reader, headers map[string]string) (*Response, error) {
	res, err := c.do(ctx, "POST", url, payload, headers)
	if err != nil {
		return nil, err
	}

	return res, nil
}

// Cached returns a new Service instance with caching enabled.
func (c client) Cached(ttl time.Duration, key ...string) Service {
	if c.cache.key != "" {
		key = append(key, c.cache.key)
	}

	k := strings.Join(key, "-")
	if k != "" {
		c.cache = &httpCache{
			client:     c.cache.client,
			enabled:    true,
			ttl:        ttl,
			key:        k,
			callerSkip: c.cache.callerSkip,
		}
	}

	return &c
}

// SetCache sets the cache key for the HTTP client.
func (c client) SetCache(key ...string) Service {
	k := strings.Join(key, "-")
	if k != "" {
		cache := *c.cache
		c.cache = &cache
		c.cache.key = k
	}

	return &c
}

// Skip sets the number of stack frames to skip when determining the caller function
// returning a new Service instance with skipped caller.
func (c client) Skip(size int) Service {
	cache := *c.cache
	c.cache = &cache
	c.cache.callerSkip += size

	return &c
}

func (c *client) do(ctx context.Context, method, url string, payload io.Reader, headers map[string]string) (*Response, error) {
	c.setCachePrefix(utils.GetCaller(c.cache.callerSkip))
	if c.cache.enabled {
		res, err := c.cache.client.Get(ctx, c.cache.key)
		if err == nil {
			return &Response{
				StatusCode: http.StatusOK,
				Body:       []byte(res),
			}, nil
		}
	}

	req, err := http.NewRequest(method, url, payload)
	if err != nil {
		return nil, err
	}

	for key, value := range headers {
		req.Header.Set(key, value)
	}

	res, err := c.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	if c.cache.enabled {
		_ = c.cache.client.Set(ctx, c.cache.ttl, c.cache.key, res.Body)
	}

	return &Response{
		StatusCode: res.StatusCode,
		Body:       body,
	}, nil
}

func (c *client) setCachePrefix(prefix string) {
	c.cache.key = fmt.Sprintf("%s-%s", prefix, c.cache.key)
}
