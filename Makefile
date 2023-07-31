.PHONY: all

all: logger_proto logger_swagger

logger_proto:
	protoc -I/usr/local/include -I./proto -I${GRPC_GATEWAY}/third_party/googleapis  \
      -I${GRPC_GATEWAY}  \
      --go-grpc_out=./proto --go-grpc_opt=paths=source_relative \
      --go_opt=paths=source_relative \
      --go_out=./proto ./proto/*.proto


logger_swagger:
	protoc -I/usr/local/include -I./proto -I${GRPC_GATEWAY}/third_party/googleapis  \
      -I${GRPC_GATEWAY} \
      --swagger_out=version=false,json_names_for_fields=false,allow_delete_body=true,merge_file_name=logger,allow_merge=true:./swagger \
      ./proto/*.proto