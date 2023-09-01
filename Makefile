.PHONY: generated-code
generated-code:
	protoc --go_out=pkg/api/ --go-grpc_out=pkg/api/ api/*.proto
