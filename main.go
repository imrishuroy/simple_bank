package main

import (
	"context"
	"database/sql"
	"log"
	"net"
	"net/http"

	_ "github.com/golang-migrate/migrate/v4/database/postgres" // required for db migration
	_ "github.com/golang-migrate/migrate/v4/source/file"       // required for db migration

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	db "github.com/imrishuroy/simplebank/db/sqlc"
	_ "github.com/imrishuroy/simplebank/doc/statik"
	"github.com/imrishuroy/simplebank/gapi"
	"github.com/imrishuroy/simplebank/pb"
	"github.com/imrishuroy/simplebank/util"
	_ "github.com/lib/pq"
	"github.com/rakyll/statik/fs"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"google.golang.org/protobuf/encoding/protojson"
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

	runDBMigration(config.MigrationURL, config.DBSource)

	store := db.NewStore(conn)

	go runGatewayServer(config, store)

	// if you want to run HTTP server
	// runHTTPServer(config, store)

	// if you want to run gRPC server
	runGrpcServer(config, store)

}

func runDBMigration(migrationURL string, dbSource string) {
	// migration, err := migrate.New(migrationURL, dbSource)
	// if err != nil {
	// 	log.Fatal("cannot create new migrate instance:", err)
	// }

	// if err = migration.Up(); err != nil && err != migrate.ErrNoChange {
	// 	log.Fatal("failed to run migrate up:", err)
	// }

	log.Println("migration completed successfully")

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
		log.Fatal("cannot create listener:", err)
	}

	log.Printf("start gRPC server at %s", listener.Addr().String())
	err = grpcServer.Serve(listener)
	if err != nil {
		log.Fatal("cannot start gRPC server:", err)
	}
}

func runGatewayServer(config util.Config, store db.Store) {
	server, err := gapi.NewServer(config, store)
	if err != nil {
		log.Fatal("cannot create server:", err)
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

	err = pb.RegisterSimpleBankHandlerServer(ctx, grpcMux, server)
	if err != nil {
		log.Fatal("cannot register handler server:", err)
	}

	mux := http.NewServeMux()
	mux.Handle("/", grpcMux)

	// fs := http.FileServer(http.Dir("./doc/swagger"))
	// mux.Handle("/swagger/", http.StripPrefix("/swagger/", fs))

	statikFS, err := fs.New()
	if err != nil {
		log.Fatal("cannot create statik fs:", err)
	}
	swaggerHandler := http.StripPrefix("/swagger", http.FileServer(statikFS))
	mux.Handle("/swagger/", swaggerHandler)

	listener, err := net.Listen("tcp", config.HTTPServerAddress)
	if err != nil {
		log.Fatal("cannot create listener:", err)
	}

	log.Printf("start HTTP gateway server at %s", listener.Addr().String())
	err = http.Serve(listener, mux)
	if err != nil {
		log.Fatal("cannot start HTTP gateway server:", err)
	}
}
