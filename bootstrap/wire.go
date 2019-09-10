//+build wireinject

package bootstrap

import (
	"github.com/google/wire"
	"github.com/stainour/test10/api"
	"github.com/stainour/test10/domain"
	"github.com/stainour/test10/infrastructure"
)

func InitializeService() (api.UrlShortenerServer, error) {
	wire.Build(infrastructure.NewMongoSequenceGenerator, infrastructure.NewMongoUrlMappingRepository, NewMongoConnectionSetting, domain.NewApiService, api.NewServer)
	return nil, nil
}

func NewMongoConnectionSetting() infrastructure.MongoConnectionSetting {
	return *infrastructure.NewMongoConnectionSetting("mongodb://localhost:27017", "urlShortener")
}
