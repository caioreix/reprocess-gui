package reprocess

import (
	"context"
	"encoding/json"
	"fmt"
	"net/url"
	"reflect"
	"time"

	"reprocess-gui/internal/errors"
	"reprocess-gui/internal/utils"
)

// ConsumerInterface defines the interface for interacting with consumers.
type ConsumerInterface interface {
	GetPagedConsumers(ctx context.Context, token string, limit int) (*PagedConsumer, error)
	GetAllPages(ctx context.Context, limit int) (consumers []*Consumer, lastToken string, err error)
}

// Consumer represents a consumer entity.
type Consumer struct {
	ID         string `json:"id"`
	Name       string `json:"name"`
	Type       string `json:"type"`
	Table      string `json:"table"`
	Team       string `json:"team"`
	Target     string `json:"target"`
	TargetType string `json:"target_type"`
	ErrorPath  string `json:"error_path"`
	RowTTL     int    `json:"row_ttl"`

	CreatedAT time.Time `json:"created_at"`
	UpdatedAT time.Time `json:"updated_at"`
}

// PagedConsumer represents a paginated list of consumers.
type PagedConsumer struct {
	Consumers  []*Consumer       `json:"consumers"`
	Pagination *utils.Pagination `json:"pagination"`
}

// GetPagedConsumers retrieves consumers in a paginated manner.
func (a *api) GetPagedConsumers(ctx context.Context, token string, limit int) (*PagedConsumer, error) {
	urlPath := a.baseURL.JoinPath("/consumers")
	values := url.Values{
		"page_token": []string{token},
		"limit":      []string{fmt.Sprint(limit)},
	}
	urlPath.RawQuery = values.Encode()

	res, err := a.httpClient.Get(ctx, urlPath.String(), nil, nil)
	if err != nil {
		return nil, err
	}
	if res.StatusCode >= 400 {
		return nil, errors.ErrStatusCode
	}

	consumers := &PagedConsumer{}
	if err := json.Unmarshal(res.Body, &consumers); err != nil {
		return nil, err
	}

	if reflect.DeepEqual(*consumers, PagedConsumer{}) {
		return nil, errors.ErrEmptyResponse
	}

	return consumers, nil
}

// GetAllPages retrieves all consumers across pages.
func (a *api) GetAllPages(ctx context.Context, limit int) (consumers []*Consumer, lastToken string, err error) {
	api := *a
	api.httpClient = api.httpClient.Skip(1)
	consumers = make([]*Consumer, 0)

	for {
		api.httpClient = api.httpClient.SetCache(fmt.Sprint(limit), lastToken)
		c, err := api.GetPagedConsumers(ctx, lastToken, limit)
		if err != nil {
			return consumers, lastToken, err
		}
		lastToken = c.Pagination.NextPage
		consumers = append(consumers, c.Consumers...)

		if lastToken == "" {
			break
		}
	}

	return consumers, "", nil
}
