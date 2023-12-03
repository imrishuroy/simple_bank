To Run sqlc ( with docker )

    docker pull kjconroy/sqlc

    docker run --rm -v ${PWD}:/src -w /src kjconroy/sqlc generate

To run gRPC client

   docker run --rm -v "$(pwd):/mount:ro" ghcr.io/ktr0731/evans:latest --path ./proto/files --proto service_simple_bank.proto --host localhost --port 9090 repl



docker run --name simplebank --network bank-network -p 8080:8080 -e GIN_MODE=release -e DB_SOURCE="postgresql://root:Prince2024@postgres12:5432/simple_bank?sslmode=disable" simplebank:latest

 To get external API
 kubectl get services simple-bank-api-service

 proto:
	rm -f pb/*.go
	Remove-Item pb/*.go -Force
	protoc --proto_path=proto --go_out=pb --go_opt=paths=source_relative \
    --go-grpc_out=pb --go-grpc_opt=paths=source_relative \
    --openapiv2_out=doc/swagger \
    proto/*.proto


    https://medium.com/@rafaellevissa/microservices-architecture-in-golang-a-transformative-experience-596d1427a29a


