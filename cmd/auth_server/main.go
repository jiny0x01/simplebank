package main

import (
	"database/sql"
	"log"

	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jiny0x01/simplebank/auth"
	db "github.com/jiny0x01/simplebank/db/sqlc"
	"github.com/jiny0x01/simplebank/util"
	_ "github.com/lib/pq"
)

func main() {
	config, err := util.LoadConfig(".")
	if err != nil {
		log.Fatal("cannot load config:", err)
	}
	conn, err := sql.Open(config.DBDriver, config.DBSource)
	if err != nil {
		log.Fatal("cannot testDB connect to db:", err)
	}
	store := db.NewStore(conn)
	// grpc_gateway 플러그인을 사용하여 http request를 grpc로 변환해주고 응답은 다시 http로 받는 식으로 구현
	//	go runGatewayServer(config, store)
	runGinServer(config, store)
}

/*
func runGatewayServer(config util.Config, store db.Store) {
	server, err := auth.NewServer(config, store)
	if err != nil {
		log.Fatal("cannot create HTTP gateway auth server: %w", err)
	}

	jsonOption := runtime.WithMarshalerOption(runtime.MIMEWildcard, &runtime.JSONPb{
		MarshalOptions: protojson.MarshalOptions{
			UseProtoNames: true,
		},
		UnmarshalOptions: protojson.UnmarshalOptions{
			DiscardUnknown: true,
		},
	})

	grpcMux := runtime.NewServeMux(jsonOption)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	err = pb.RegisterAuthHandlerServer(ctx, grpcMux, server)
	if err != nil {
		log.Fatal("cannot register handler HTTP gateway auth server: %w", err)
	}

	mux := http.NewServeMux()
	mux.Handle("/", grpcMux)

	fs := http.FileServer(http.Dir("./doc/swagger"))
	mux.Handle("/swagger/", http.StripPrefix("/swagger/", fs))

	listener, err := net.Listen("tcp", config.GRPCAuthServerAddress)
	if err != nil {
		log.Fatal("cannot create listener: %w", err)
	}

	log.Printf("start HTTP gateway server at %s", listener.Addr().String())
	err = http.Serve(listener, mux)
	if err != nil {
		log.Fatal("cannot start HTTP gateway auth server: %w", err)
	}
}
*/

func runGinServer(config util.Config, store db.Store) {
	server, err := auth.NewServer(config, store)
	if err != nil {
		log.Fatal("cannot create auth server:", err)
	}

	err = server.Start(config.HTTPAuthServerAddress)
	if err != nil {
		log.Fatal("cannot start auth server:", err)
	}
}
