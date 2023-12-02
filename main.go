package main

import (
	"database/sql"
	"log"
	"net"

	db "github.com/imrishuroy/simplebank/db/sqlc"
	"github.com/imrishuroy/simplebank/gapi"
	"github.com/imrishuroy/simplebank/pb"
	"github.com/imrishuroy/simplebank/util"
	_ "github.com/lib/pq"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

// const (
// 	dbDriver      = "postgres"
// 	dbSource      = "postgresql://root:Prince2024@localhost:5432/simple_bank?sslmode=disable"
// 	serverAddress = "localhost:8080"
// )

func main() {

	config, err := util.LoadConfig(".") // . means current directory

	if err != nil {
		log.Fatal("cannot load config:", err)
	}

	conn, err := sql.Open(config.DBDriver, config.DBSource)
	if err != nil {
		log.Fatal("cannot connect to db:", err)
	}

	store := db.NewStore(conn)

	// if you want to run HTTP server
	// runHTTPServer(config, store)

	// if you want to run gRPC server
	runGrpcServer(config, store)

}

// func runHTTPServer(config util.Config, store db.Store) {

// 	server, err := api.NewServer(config, store)
// 	if err != nil {
// 		log.Fatal("cannot create server:", err)
// 	}

// 	err = server.Start(config.HTTPServerAddress)
// 	if err != nil {
// 		log.Fatal("cannot start server:", err)
// 	}

// }

func runGrpcServer(config util.Config, store db.Store) {
	server, err := gapi.NewServer(config, store)
	if err != nil {
		log.Fatal("cannot create server:", err)
	}

	grpcServer := grpc.NewServer()
	pb.RegisterSimpleBankServer(grpcServer, server)
	reflection.Register(grpcServer)

	listener, err := net.Listen("tcp", config.GRPCServerAddress)
	if err != nil {
		log.Fatal("cannot create listener")
	}

	log.Printf("start gRPC server at %s", listener.Addr().String())
	err = grpcServer.Serve(listener)
	if err != nil {
		log.Fatal("cannot start gRPC server")
	}
}
