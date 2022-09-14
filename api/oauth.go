package api

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jiny0x01/simplebank/auth"
	db "github.com/jiny0x01/simplebank/db/sqlc"
	"golang.org/x/oauth2"
)

func (server *Server) renderAuthView(ctx *gin.Context) {
	// Redirect user to consent page to ask for permission
	// for the scopes specified above.
	url := auth.Conf.AuthCodeURL("state", oauth2.AccessTypeOffline)
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

	token, err := auth.Authenticate(code)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	client := auth.Conf.Client(ctx, token)
	// client.Get("api url")
	// scope에 동의한 정보면 다 갖고올 수 있다.
	userInfoResp, err := client.Get(auth.UserInfoAPIEndpoint)
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
	server.renderAuthView(ctx)
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
	id_token, err := auth.Verify(req.AccessToken)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	sub := id_token["sub"].(string)
	fmt.Printf("id:%s\n", sub)

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

func (server *Server) refreshAccessToken(ctx *gin.Context) {
	refresh_token := ctx.Request.FormValue("refresh_token")
	fmt.Printf("refresh:%s\n", refresh_token)
	token, err := auth.Refresh(refresh_token)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	ctx.JSON(http.StatusOK, token)
}

func (server *Server) revokeAccessToken(ctx *gin.Context) {
	token := ctx.Request.FormValue("token")
	if token == "" {
		ctx.JSON(http.StatusBadRequest, "FromValue require token which is access token or refresh token both ok")
		return
	}
	err := auth.Revoke(auth.RevokeGoogleAPIEndpoint, token)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, nil)
}
