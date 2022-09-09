package api

import (
	"fmt"

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

func (server *Server) setupRouter() {
	router := gin.Default()
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		v.RegisterValidation("currency", validCurrency)
	}

	router.POST("/users", server.createUser)
	router.POST("/users/login", server.loginUser)

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
