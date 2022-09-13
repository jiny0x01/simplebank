package auth

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"net/url"

	"github.com/spf13/viper"
	"golang.org/x/oauth2"
)

type OauthEnv struct {
	ClientID     string `mapstructure:"CLIENT_ID"`
	ClientSecret string `mapstructure:"CLIENT_SECRET"`
}

const (
	redirectURL = "http://localhost:8080/users/auth/callback"

	// 인증 후 유저 정보를 가져오기 위한 API
	userInfoAPIEndpoint = "https://www.googleapis.com/oauth2/v3/userinfo"

	// 인증 권한 범위. 여기에서는 프로필 정보 권한만 사용
	scopeEmail              = "https://www.googleapis.com/auth/userinfo.email"
	scopeProfile            = "https://www.googleapis.com/auth/userinfo.profile"
	authURL                 = "https://accounts.google.com/o/oauth2/auth"
	tokenURL                = "https://oauth2.googleapis.com/token" // oauth2.Config.Exchange에서 내부적으로 사용함
	revokeGoogleAPIEndpoint = "https://oauth2.googleapis.com/revoke"
)

var Conf oauth2.Config

func init() {
	viper.AddConfigPath(".")
	viper.SetConfigName("app")
	viper.SetConfigType("env")
	viper.AutomaticEnv()
	err := viper.ReadInConfig()
	if err != nil {
		panic(err)
	}
	env := &OauthEnv{}
	err = viper.Unmarshal(&env)
	if err != nil {
		panic(err)
	}

	Conf = oauth2.Config{
		ClientID:     env.ClientID,
		ClientSecret: env.ClientSecret,
		Scopes: []string{
			scopeEmail,
			scopeProfile,
		},
		Endpoint: oauth2.Endpoint{
			AuthURL:  authURL,
			TokenURL: tokenURL,
		},
		RedirectURL: redirectURL,
	}
}

func Authenticate(code string) (*oauth2.Token, error) {
	if code == "" {
		return nil, errors.New("code is not exist")
	}
	fmt.Println(code)
	// Exchange()의 역할은 아래 url 단계에 해당함
	// https://developers.google.com/identity/protocols/oauth2/web-server#exchange-authorization-code
	token, err := Conf.Exchange(context.Background(), code)
	if err != nil {
		return nil, err
	}
	fmt.Printf("oauth token\n\taccess:%s\n\trefresh:%s\n", token.AccessToken, token.RefreshToken)
	// token은 redis나 rdb에 저장하는게 맞다.
	// access token은 client에게(웹브라우저, 모바일 디바이스)에 보관해도 되지만
	// access token은 key-value cache db(redis)에 저장 관리해도 되지만
	// refresh token은 expiry가 기므로 rdb에 저장하는게 나아보인다.
	return token, nil
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
