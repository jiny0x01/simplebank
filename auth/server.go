package auth

import (
	"fmt"
	"text/template"

	"github.com/gin-gonic/gin"
	db "github.com/jiny0x01/simplebank/db/sqlc"
	"github.com/jiny0x01/simplebank/token"
	"github.com/jiny0x01/simplebank/util"
	"golang.org/x/oauth2"
)

type Server struct {
	config      util.Config
	oauthConfig oauth2.Config
	store       db.Store
	tokenMaker  token.Maker
	router      *gin.Engine
}

func (server *Server) RenderTemplate(ctx *gin.Context, name string, data interface{}) {
	tmpl, _ := template.ParseFiles(name)
	tmpl.Execute(ctx.Writer, data)
}

func NewServer(config util.Config, store db.Store) (*Server, error) {
	tokenMaker, err := token.NewJWTMaker(config.TokenSymmetricKey)
	if err != nil {
		return nil, fmt.Errorf("cannot create token maker: %w", err)
	}
	server := &Server{
		config: config,
		oauthConfig: oauth2.Config{
			ClientID:     config.OauthClientID,
			ClientSecret: config.OauthClientSecret,
			Scopes: []string{
				scopeEmail,
				scopeProfile,
			},
			Endpoint: oauth2.Endpoint{
				AuthURL:  authEndpoint,
				TokenURL: tokenEndpoint,
			},
			RedirectURL: redirectURL,
		},
		store:      store,
		tokenMaker: tokenMaker,
	}
	server.setupRouter()
	return server, nil
}

func (server *Server) setupRouter() {
	router := gin.Default()

	router.GET("/", server.renderOauthView)
	router.GET("/oauth/callback", server.callbackOauth)
	router.GET("/users/oauth/info", server.getOauthUserInfo)
	router.POST("/tokens/oauth/revoke", server.revokeOauthAccessToken)
	router.POST("/tokens/oauth/refresh", server.refreshOauthAccessToken)
	router.POST("/tokens/renew_access/", server.renewAccessToken)
	router.GET("/users/login/auth", server.loginOauthUser)
	server.router = router
}

func (server *Server) Start(address string) error {
	return server.router.Run(address)
}

func errorResponse(err error) gin.H {
	return gin.H{"error": err.Error()}
}
