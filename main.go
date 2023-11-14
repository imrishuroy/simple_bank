package main

import (
	"database/sql"
	"log"

	"github.com/imrishuroy/simplebank/api"
	db "github.com/imrishuroy/simplebank/db/sqlc"
	_ "github.com/lib/pq"
)

const (
	dbDriver      = "postgres"
	dbSource      = "postgresql://root:Prince2024@localhost:5432/simple_bank?sslmode=disable"
	serverAddress = "localhost:8080"
)

func main() {

	conn, err := sql.Open(dbDriver, dbSource)
	if err != nil {
		log.Fatal("cannot connect to db:", err)
	}

	store := db.NewStore(conn)
	server := api.NewServer(store)

	err = server.Start(serverAddress)
	if err != nil {
		log.Fatal("cannot start server:", err)
	}
}
