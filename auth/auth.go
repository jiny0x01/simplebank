package auth

import (
	"errors"
	"fmt"
	"net/http"
	"net/url"
)

func Authenticate() (err error) {
	return nil
}

func Revoke(endpoint string, token string) error {
	res, err := http.PostForm(endpoint, url.Values{"token": {token}})
	if err != nil {
		return err
	}
	defer res.Body.Close()

	fmt.Printf("status:%d\n", res.StatusCode)
	if res.StatusCode != 200 {
		return errors.New("invalid_token")
	}
	return nil
}
