To Run sqlc ( with docker )

    docker pull kjconroy/sqlc

    docker run --rm -v ${PWD}:/src -w /src kjconroy/sqlc generate

To run gRPC client

   docker run --rm -v "$(pwd):/mount:ro" ghcr.io/ktr0731/evans:latest --path ./proto/files --proto service_simple_bank.proto --host localhost --port 9090 repl

Test redis server

    docker exec -it redis redis-cli ping
