package api

import (
	"fmt"
	"html/template"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
	db "github.com/jiny0x01/simplebank/db/sqlc"
	"github.com/jiny0x01/simplebank/token"
	"github.com/jiny0x01/simplebank/util"
)

type Server struct {
	config     util.Config
	store      db.Store
	tokenMaker token.Maker
	router     *gin.Engine
}

func NewServer(config util.Config, store db.Store) (*Server, error) {
	tokenMaker, err := token.NewJWTMaker(config.TokenSymmetricKey) //env
	if err != nil {
		return nil, fmt.Errorf("cannot create token maker: %w", err)
	}
	server := &Server{
		config:     config,
		store:      store,
		tokenMaker: tokenMaker,
	}
	server.setupRouter()
	return server, nil
}

func (server *Server) RenderTemplate(ctx *gin.Context, name string, data interface{}) {
	tmpl, _ := template.ParseFiles(name)
	tmpl.Execute(ctx.Writer, data)
}

func (server *Server) RenderMainView(ctx *gin.Context) {
	server.RenderTemplate(ctx, "html/main.html", nil)
}

func (server *Server) setupRouter() {
	router := gin.Default()
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		v.RegisterValidation("currency", validCurrency)
	}
	// 인증 기능들 분리
	// 우선은 redirect로 구현하고 후에 리펙토링 할 때 reverse proxy가 필요한 부분은 수정
	router.POST("/users", func(ctx *gin.Context) {
		// 307(StatusTemporaryRedirect) will
		// reissue the same request and verb to a different URI specified by the location header.
		// API 서버로 회원가입/로그인 요청을 인증서버로 redirect
		ctx.Redirect(http.StatusTemporaryRedirect, 
			server.config.HTTPAuthServerAddress + "/v1/create_user")
	})
	router.POST("/users/login", func(ctx *gin.Context) {
		ctx.Redirect(http.StatusTemporaryRedirect, 
			server.config.HTTPAuthServerAddress + "/v1/login_user")
	})

	router.GET("/users/auth", func(ctx *gin.Context) {
		ctx.Redirect(http.StatusTemporaryRedirect, 
			server.config.HTTPAuthServerAddress + "/users/auth")
	})

	/*
		router.GET("/users/auth/refresh", server.refreshAccessToken)
		router.GET("/users/auth/callback", server.callbackOauth)
		router.GET("/users/auth/revoke", server.revokeAccessToken)
		router.POST("/users/auth/", server.getOauthUserInfo)
		router.GET("/users/login/auth", server.loginOauthUser)
	*/
	router.POST("/tokens/renew_access/", func(ctx *gin.Context) {
		ctx.Redirect(http.StatusTemporaryRedirect, 
			server.config.HTTPAuthServerAddress + "/tokens/renew_access")
	})

	router.GET("/", server.RenderMainView)
	authRoutes := router.Group("/").Use(authMiddleware(server.tokenMaker))

	authRoutes.POST("/accounts", server.createAccount)
	authRoutes.GET("/accounts/:id", server.getAccount)
	authRoutes.GET("/accounts", server.listAccount)
	// PATCH는 이미 존재하는 리소스의 일부를 바꿀 때만 사용.
	// PUT은 payload에 있는 값으로 replace하는 용도.
	// PUT으로 요청한 URI에 리소스가(데이터) 없으면 생성하고 201(Created 응답)
	// PUT으로 요청한 URI에 리소스가 있으면 payload에 담긴 값으로 변경하고 200(ok)나 204(no content) 이용
	authRoutes.PATCH("/accounts", server.updateAccount)
	authRoutes.DELETE("/accounts", server.deleteAccount)

	router.POST("/transfers", server.createTransfer)
	server.router = router

}

func (server *Server) Start(address string) error {
	return server.router.Run(address)
}

func errorResponse(err error) gin.H {
	return gin.H{"error": err.Error()}
}
