package gapi

import (
	"fmt"

	db "github.com/jiny0x01/simplebank/db/sqlc"
	"github.com/jiny0x01/simplebank/pb"
	"github.com/jiny0x01/simplebank/token"
	"github.com/jiny0x01/simplebank/util"
)

type Server struct {
	pb.UnimplementedSimpleBankServer
	config     util.Config
	store      db.Store
	tokenMaker token.Maker
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

	return server, nil
}
