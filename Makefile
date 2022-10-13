.PHONY: test

test:
	go test ./...

cover:
	go test -cover ./...

proto:
	protoc --go_out=. --go_opt=paths=source_relative \
	  --go-grpc_out=. --go-grpc_opt=paths=source_relative \
	  proto/person.proto proto/user.proto

