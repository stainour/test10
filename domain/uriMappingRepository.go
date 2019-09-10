package domain

import (
	"context"
	"net/url"
)

const (
	OK AddResult = iota
	AlreadyExists
	Fail
)

type AddResult int

type UrlMappingRepository interface {
	AddIfNotExists(context context.Context, uriMapping *UrlMapping) (AddResult, error)
	FindById(context context.Context, id *url.URL) (*UrlMapping, error)
	IncrementHitCount(context context.Context, id *url.URL) error
	GetAll(context context.Context) (values <-chan *UrlMapping, errors <-chan error)
	FindByShortenedKey(context context.Context, key string) (*UrlMapping, error)
}
