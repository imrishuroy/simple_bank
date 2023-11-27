To Run sqlc ( with docker )

    docker pull kjconroy/sqlc

    docker run --rm -v ${PWD}:/src -w /src kjconroy/sqlc generate



docker run --name simplebank --network bank-network -p 8080:8080 -e GIN_MODE=release -e DB_SOURCE="postgresql://root:Prince2024@postgres12:5432/simple_bank?sslmode=disable" simplebank:latest