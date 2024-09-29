package client

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
)

type Page[T any] interface {
	GetNextToken() string
	GetData() []T
}

type PaginatorOptions struct {
	Sort     string
	PageSize int
}

func NewPaginator[T any, P Page[T]](client Client, endpoint string, options PaginatorOptions) *Paginator[T, P] {
	return &Paginator[T, P]{
		client:   client,
		endpoint: endpoint,
		options:  options,
		first:    true,
		token:    "",
	}
}

type Paginator[T any, P Page[T]] struct {
	client   Client
	endpoint string
	first    bool
	token    string
	options  PaginatorOptions
}

func (p *Paginator[T, P]) HasNext() bool {
	return p.first || p.token != ""
}

func (p *Paginator[T, P]) Next(ctx context.Context) ([]T, error) {
	p.first = false

	query := url.Values{}

	if p.options.Sort != "" {
		query.Set("sort", string(p.options.Sort))
	}

	if p.options.PageSize > 0 {
		query.Set("page_size", fmt.Sprint(p.options.PageSize))
	}

	if p.token != "" {
		query.Set("next_token", p.token)
	}

	req, err := p.client.NewRequest(ctx, http.MethodGet, p.endpoint, query, nil)
	if err != nil {
		return nil, err
	}

	var response P
	if err := p.client.DoRequest(req, &response); err != nil {
		return nil, err
	}

	p.token = response.GetNextToken()
	data := response.GetData()
	if data == nil {
		return []T{}, nil
	}

	return data, nil
}
