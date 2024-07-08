package utils

import (
	"net/http"
	"strconv"

	customerr "github.com/berkayaydmr/language-learning-api/pkg/error"
)

type UrlParamKey string

const (
	UrlParamKeyID UrlParamKey = "id"
)

func (k UrlParamKey) String() string {
	return string(k)
}

func ParseStrToInt(str string) (int, error) {
	return strconv.Atoi(str)
}

func GetUrlParam(r *http.Request, key UrlParamKey) (*string, error) {
	value := r.FormValue(key.String())
	if value == "" {
		return nil, customerr.ErrInvalidParameter
	}

	return &value, nil
}
