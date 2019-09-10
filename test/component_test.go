package test

import (
	"context"
	"github.com/stainour/test10/api"
	"github.com/stretchr/testify/assert"
	"io"
	"net/url"
	"testing"
	"time"
)

var testUrl *url.URL
var grpcClient api.UrlShortenerClient

func init() {
	testUrl, _ = url.Parse("https://www.google.com")

	go StartServer()
	for {
		client, e := GetClient()
		if e != nil {
			time.Sleep(time.Second)
			continue
		}
		grpcClient = client
		break
	}
}

func TestApi(t *testing.T) {
	context := context.Background()
	_, err := grpcClient.Create(context, &api.OriginalUrlRequest{
		OriginalUrl: testUrl.String(),
	})

	if err != nil {
		t.Error(err)
	}
	stream, err := grpcClient.GetStatistic(context, &api.StatisticRequest{})
	count := 0
	for {
		stat, err := stream.Recv()

		if err == io.EOF {
			break
		}
		if err != nil && err != io.EOF {
			t.Error(err)
		}
		assert.NotNil(t, stat)
		assert.NotEmpty(t, stat.Url)
		count++
	}
	assert.Greater(t, count, 0)

	_, err = grpcClient.GetOriginalUrl(context, &api.ShortenedlUrlRequest{
		ShortenedlUrl: "url",
	})

	if err != nil {
		t.Error(err)
	}
}
