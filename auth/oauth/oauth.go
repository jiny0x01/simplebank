package auth

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"

	"github.com/jiny0x01/simplebank/util"
	"golang.org/x/oauth2"
)
func init() {
	// FIXME: util.LoadConfig(".") -> make test 할 때 panic 발생
	// util.LoadConfig("../") -> make server 할 때 panic 발생
	// server.config로 가져오는 방식이나 다른 방식 생각 필요
	envConf, err := util.LoadConfig(".")

	if err != nil {
		panic(err)
	}

	Conf = oauth2.Config{
		ClientID:     envConf.ClientID,
		ClientSecret: envConf.ClientSecret,
		Scopes: []string{
			scopeEmail,
			scopeProfile,
		},
		Endpoint: oauth2.Endpoint{
			AuthURL:  authEndpoint,
			TokenURL: tokenEndpoint,
		},
		RedirectURL: redirectURL,
	}
}
