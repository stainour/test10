package domain

import (
	"context"
	"errors"
	"net/url"
)

type CreateResult int

const (
	Conflict CreateResult = iota
	Inserted
	Error
)

type ApiService struct {
	sequenceGenerator SequenceNumberGenerator
	repository        UrlMappingRepository
}

func (service *ApiService) ResolveByKey(context context.Context, shortenedKey string) (*UrlMapping, error) {
	if shortenedKey == "" {
		return nil, errors.New("shortenedKey is empty")
	}
	mapping, err := service.repository.FindByShortenedKey(context, shortenedKey)

	if err != nil {
		return nil, err
	}

	if mapping == nil {
		return nil, nil
	}
	url, err := url.Parse(mapping.Uri())

	if err != nil {
		return nil, err
	}

	err = service.repository.IncrementHitCount(context, url)
	if err != nil {
		return nil, err
	}

	return mapping, nil
}

func (service *ApiService) Create(context context.Context, url *url.URL) (result CreateResult, error error) {
	if url == nil {
		return Error, errors.New("url is nil")
	}

	mapping, err := service.repository.FindById(context, url)

	if err != nil {
		return Error, err
	}

	if mapping != nil {
		return Conflict, nil
	}

	nextValue, err := service.sequenceGenerator.NextValue(context)
	if err != nil {
		return Error, err
	}

	urlMapping := NewUrlMapping(url, nextValue)

	addResult, err := service.repository.AddIfNotExists(context, urlMapping)
	if err != nil {
		return Error, nil
	}

	if addResult == OK {
		return Inserted, nil
	}

	return Conflict, nil
}

func (service *ApiService) GetAll(context context.Context) (values <-chan *UrlMapping, errors <-chan error) {
	return service.repository.GetAll(context)
}

func NewApiService(sequenceGenerator SequenceNumberGenerator, mappingRepository UrlMappingRepository) *ApiService {
	return &ApiService{
		sequenceGenerator: sequenceGenerator,
		repository:        mappingRepository,
	}
}
