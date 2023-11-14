package main

import (
	"database/sql"
	"log"

	"github.com/imrishuroy/simplebank/api"
	db "github.com/imrishuroy/simplebank/db/sqlc"
	"github.com/imrishuroy/simplebank/util"
	_ "github.com/lib/pq"
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
	server := api.NewServer(store)

	err = server.Start(config.ServerAddress)
	if err != nil {
		log.Fatal("cannot start server:", err)
	}
}
