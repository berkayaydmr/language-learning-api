package utils

import (
	"net/http"

	customerr "github.com/berkayaydmr/language-learning-api/pkg/error"
)

type UrlParamKey string

const (
	UrlParamKeyID UrlParamKey = "id"
)

func (k UrlParamKey) String() string {
	return string(k)
}

func GetUrlParam(r *http.Request, key UrlParamKey) (string, error) {
	value := r.PathValue(key.String())
	if value == "" {
		return "", customerr.ErrInvalidParameter
	}

	return value, nil
}
