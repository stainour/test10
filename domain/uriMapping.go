package domain

import (
	"net/url"
	"strings"
)

type UrlMapping struct {
	hitCount     int64
	shortenedKey string
	uri          string
}

func NewUrlMapping(url *url.URL, uniqueNumber int64) *UrlMapping {
	convertToBase62 := func(value int64) string {
		const alphabet = "i8kXgbjRKvEh6M7UsVaSdAJpw59cuZnBLrPNoDzmfHxIG1lYCyFQ23qWOe0T4t"
		const base int64 = 62

		stringBuilder := strings.Builder{}
		for value > 0 {
			stringBuilder.WriteString(string(alphabet[value%base]))
			value /= base
		}
		return stringBuilder.String()
	}

	return &UrlMapping{uri: url.String(), shortenedKey: convertToBase62(uniqueNumber)}
}
func NewFullUrlMapping(url *url.URL, shortenedKey string, hitCount int64) *UrlMapping {
	return &UrlMapping{uri: url.String(), shortenedKey: shortenedKey, hitCount: hitCount}
}

func (u *UrlMapping) HitCount() int64 {
	return u.hitCount
}

func (u *UrlMapping) Uri() string {
	return u.uri
}

func (u *UrlMapping) ShortenedKey() string {
	return u.shortenedKey
}
