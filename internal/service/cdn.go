package service

import (
	"fmt"
	"net/http"
)

type CDNService struct {
	base string
}

func NewCDNService(base string) *CDNService {
	return &CDNService{base: base}
}

func (s *CDNService) AssetURL(objectKey string) string {
	return fmt.Sprintf("%s/assets/%s", s.base, objectKey)
}

func (s *CDNService) SetCacheHeaders(w http.ResponseWriter, kind string) {
	switch kind {
	case "art":
		w.Header().Set("Cache-Control", "public, max-age=604800, stale-while-revalidate=86400")
	case "meta":
		w.Header().Set("Cache-Control", "public, max-age=60, stale-while-revalidate=10")
	default:
		w.Header().Set("Cache-Control", "no-store")
	}
	w.Header().Set("Vary", "Accept-Encoding")
}
