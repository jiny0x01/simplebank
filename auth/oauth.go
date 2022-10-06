package auth

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"

	"github.com/gin-gonic/gin"
	db "github.com/jiny0x01/simplebank/db/sqlc"
	"golang.org/x/oauth2"
)

const (
	redirectURL = "http://localhost:8081/oauth/callback"
	// 인증 권한 범위. 여기에서는 프로필 정보 권한만 사용
	scopeEmail   = "https://www.googleapis.com/auth/userinfo.email"
	scopeProfile = "https://www.googleapis.com/auth/userinfo.profile"

	authEndpoint            = "https://accounts.google.com/o/oauth2/auth"
	tokenEndpoint           = "https://oauth2.googleapis.com/token"     // oauth2.Config.Exchange에서 내부적으로 사용함
	tokenInfoEndpoint       = "https://oauth2.googleapis.com/tokeninfo" // ?access_token="accesstoken"
	RevokeGoogleAPIEndpoint = "https://oauth2.googleapis.com/revoke"

	// 인증 후 유저 정보를 가져오기 위한 API
	UserInfoAPIEndpoint = "https://www.googleapis.com/oauth2/v3/userinfo"
)

func (server *Server) renderOauthView(ctx *gin.Context) {
	// Redirect user to consent page to ask for permission
	// for the scopes specified above.
	url := server.oauthConfig.AuthCodeURL("state", oauth2.AccessTypeOffline)
	// Online is the default if neither is specified.
	// If your application needs to refresh access tokens
	// when the user is not present at the browser, then use offline.
	// This will result in your application obtaining a refresh token
	// the first time your application exchanges an authorization code for a user.

	// access token이 만료되면 AccessType이 offline이면 사용자가 없어도 access token을 새로 갱신 할 수 있다.
	fmt.Printf("Visit the URL for the auth dialog: %v", url)

	server.RenderTemplate(ctx, "html/auth.html", url)

	ctx.JSON(http.StatusOK, nil)
}

func (server *Server) callbackOauth(ctx *gin.Context) {
	code := ctx.Request.FormValue("code")

	token, err := server.Authenticate(code)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	client := server.oauthConfig.Client(ctx, token)
	// client.Get("api url")
	// scope에 동의한 정보면 다 갖고올 수 있다.
	userInfoResp, err := client.Get(UserInfoAPIEndpoint)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	defer userInfoResp.Body.Close()
	userinfo, err := ioutil.ReadAll(userInfoResp.Body)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	fmt.Printf("authUser : %s\n", string(userinfo))

	var decodeData struct {
		Sub   string `json:"sub" binding:"required"`
		Email string `json:"email" binding:"required"`
		Name  string `json:"name" binding:"required"`
	}

	err = json.Unmarshal(userinfo, &decodeData)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	// 사용자가 엑세스 권한을 취소하지 않는한 토큰 서버는 새 refresh token을 새로 반환하지 않음
	// refresh_token은 DB에 저장하고 유효하면 계속 사용하는게 좋음
	// 이미 사용자가 있으면 access_token만 지급

	authUser := db.CreateOauthUserParams{
		ID:           decodeData.Sub,
		Email:        decodeData.Email,
		Fullname:     decodeData.Name,
		Provider:     "google",
		RefreshToken: token.RefreshToken,
	}
	_, err = server.store.GetOauthUser(ctx, decodeData.Sub)
	if err != nil {
		if err == sql.ErrNoRows {
			_, err = server.store.CreateOauthUser(ctx, authUser)
			if err != nil {
				ctx.JSON(http.StatusInternalServerError, errorResponse(err))
				return
			}
		} else {
			ctx.JSON(http.StatusInternalServerError, errorResponse(err))
			return
		}
	}

	ctx.JSON(http.StatusOK, gin.H{
		"access_token": token.AccessToken,
	})
}

func (server *Server) loginOauthUser(ctx *gin.Context) {
	// 새로운 access_token을 발급한다.
	server.renderOauthView(ctx)
}

type getOauthUserInfoRequest struct {
	AccessToken string `json:"access_token" binding:"required"`
}

func (server *Server) getOauthUserInfo(ctx *gin.Context) {
	var req getOauthUserInfoRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	id_token, err := server.Verify(req.AccessToken)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	sub := id_token["sub"].(string)

	user, err := server.store.GetOauthUser(ctx, sub)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, user)
}

func (server *Server) refreshOauthAccessToken(ctx *gin.Context) {
	refresh_token := ctx.Request.FormValue("refresh_token")
	fmt.Printf("refresh:%s\n", refresh_token)
	token, err := server.Refresh(refresh_token)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	ctx.JSON(http.StatusOK, token)
}

func (server *Server) revokeOauthAccessToken(ctx *gin.Context) {
	token := ctx.Request.FormValue("token")
	if token == "" {
		ctx.JSON(http.StatusBadRequest, "FromValue require token which is access token or refresh token both ok")
		return
	}
	err := server.Revoke(RevokeGoogleAPIEndpoint, token)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, nil)
}

func (server *Server) Authenticate(code string) (*oauth2.Token, error) {
	if code == "" {
		return nil, errors.New("code is not exist")
	}
	fmt.Println(code)
	// Exchange()의 역할은 아래 url 단계에 해당함
	// https://developers.google.com/identity/protocols/oauth2/web-server#exchange-authorization-code
	token, err := server.oauthConfig.Exchange(context.Background(), code)
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

func (server *Server) Verify(access_token string) (map[string]any, error) {

	res, err := http.Get(tokenInfoEndpoint + "?access_token=" + access_token)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	var result any

	bytes, _ := ioutil.ReadAll(res.Body)
	err = json.Unmarshal(bytes, &result)
	if err != nil {
		return nil, err
	}
	fmt.Printf("Verify result:%v\n", result)
	if res.StatusCode != 200 {
		return nil, errors.New(fmt.Sprint(result))
	}
	return result.(map[string]any), nil
}

func (server *Server) Revoke(endpoint string, token string) error {
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

func (server *Server) Refresh(refreshToken string) (any, error) {
	res, err := http.PostForm(tokenEndpoint, url.Values{
		"client_id":     {server.oauthConfig.ClientID},
		"client_secret": {server.oauthConfig.ClientSecret},
		"grant_type":    {"refresh_token"},
		"refresh_token": {refreshToken},
	})
	if err != nil {
		return nil, err
	}

	defer res.Body.Close()
	fmt.Printf("status:%d\n", res.StatusCode)

	var result any

	bytes, _ := ioutil.ReadAll(res.Body)
	err = json.Unmarshal(bytes, &result)
	if err != nil {
		return nil, err
	}
	fmt.Printf("Refresh result:%v\n", result)
	if res.StatusCode != 200 {
		return nil, errors.New(fmt.Sprint(result))
	}
	return result, nil
}
