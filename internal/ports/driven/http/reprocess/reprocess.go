package reprocess

import (
	"net/url"
	"strings"
	"time"

	"reprocess-gui/internal/ports/driven/http"
)

// Service defines the interface for the reprocess service, including methods for interacting with consumers and caching.
type Service interface {
	ConsumerInterface

	Cached(key ...string) Service
	CachedTTL(ttl time.Duration, key ...string) Service
}

type api struct {
	httpClient http.Service
	cacheTTL   time.Duration
	baseURL    *url.URL
}

// NewService creates a new instance of the reprocess service.
func NewService(client http.Service, cacheTTL time.Duration, baseURL *url.URL) Service {
	return &api{
		httpClient: client,
		cacheTTL:   cacheTTL,
		baseURL:    baseURL,
	}
}

// Cached returns a new instance of the reprocess service with caching enabled.
func (a api) Cached(key ...string) Service {
	a.httpClient = a.httpClient.Cached(a.cacheTTL, strings.Join(key, "-"))
	return &a
}

// CachedTTL returns a new instance of the reprocess service with caching enabled.
func (a api) CachedTTL(ttl time.Duration, key ...string) Service {
	a.httpClient = a.httpClient.Cached(ttl, strings.Join(key, "-"))
	return &a
}
