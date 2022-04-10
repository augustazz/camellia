
build:
	CGO_ENABLED=0 go build main.go
run:
	CGO_ENABLED=0 go run main.go

gen_pb:
	protoc protocol/*.proto --go_out=.