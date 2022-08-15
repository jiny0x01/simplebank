package api

import (
	"github.com/gin-gonic/gin"
	db "github.com/jiny0x01/simplebank/db/sqlc"
)

type Server struct {
	store db.Store
	router *gin.Engine
}

func NewServer(store db.Store) *Server {
	server := &Server{store: store}
	router := gin.Default()

	router.POST("/accounts", server.createAccount)
	router.GET("/accounts/:id", server.getAccount)
	router.GET("/accounts", server.listAccount)
	// PATCH는 이미 존재하는 리소스의 일부를 바꿀 때만 사용. 
	// PUT은 payload에 있는 값으로 replace하는 용도.
	// PUT으로 요청한 URI에 리소스가(데이터) 없으면 생성하고 201(Created 응답)
	// PUT으로 요청한 URI에 리소스가 있으면 payload에 담긴 값으로 변경하고 200(ok)나 204(no content) 이용 
	router.PATCH("/accounts", server.updateAccount) 
	router.DELETE("/accounts", server.deleteAccount)

	server.router = router
	return server
}

func (server *Server) Start(address string) error {
	return server.router.Run(address)
}

func errorResponse(err error) gin.H {
	return gin.H{"error": err.Error()}
}