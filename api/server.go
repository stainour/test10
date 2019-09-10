package api

import (
	"context"
	"github.com/stainour/test10/domain"
	"net/url"
)

type urlShortenerServer struct {
	service *domain.ApiService
}

func (server *urlShortenerServer) GetStatistic(request *StatisticRequest, stream UrlShortener_GetStatisticServer) error {
	context := stream.Context()
	service := server.service
	urlMappings, errors := service.GetAll(context)

	for urlMapping := range urlMappings {
		err := stream.Send(&UrlKeyStat{
			HitCount: urlMapping.HitCount(),
			Url:      urlMapping.Uri(),
		})
		if err != nil {
			return err
		}
	}

	for err := range errors {
		return err
	}

	return nil
}

func (server *urlShortenerServer) Create(context context.Context, request *OriginalUrlRequest) (*CreationResult, error) {
	service := server.service

	url, err := url.Parse(request.OriginalUrl)
	if err != nil {
		return &CreationResult{
			Status: CreationResult_ERROR,
		}, err
	}

	createResult, err := service.Create(context, url)

	if err != nil {
		return &CreationResult{
			Status: CreationResult_ERROR,
		}, err
	}
	var result CreationResult_Status
	if createResult == domain.Inserted {
		result = CreationResult_OK
	} else {
		result = CreationResult_CONFLICT
	}

	return &CreationResult{
		Status: result,
	}, nil
}

func (server *urlShortenerServer) GetOriginalUrl(context context.Context, request *ShortenedlUrlRequest) (*OriginalUrlResponse, error) {
	mapping, err := server.service.ResolveByKey(context, request.ShortenedlUrl)
	if err != nil {
		return &OriginalUrlResponse{}, err
	}
	if mapping == nil {
		return &OriginalUrlResponse{}, nil
	}

	return &OriginalUrlResponse{
		OriginalUrl: mapping.Uri(),
	}, nil

}

func NewServer(apiService *domain.ApiService) UrlShortenerServer {
	return &urlShortenerServer{service: apiService}
}
