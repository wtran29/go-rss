package auth

import (
	"errors"
	"net/http"
	"strings"
)

func GetAPIKey(headers http.Header) (string, error) {
	authorize := headers.Get("Authorization")
	if authorize == "" {
		return "", errors.New("no authorization info found")
	}

	authorizeVals := strings.Split(authorize, " ")
	if len(authorizeVals) != 2 {
		return "", errors.New("malformed auth header")
	}
	if authorizeVals[0] != "ApiKey" {
		return "", errors.New("malformed first value of auth header")
	}
	return authorizeVals[1], nil

}
