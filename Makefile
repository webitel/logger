.PHONY: all

all: logger_proto logger_swagger

logger_proto:
	protoc -I/usr/local/include -I./proto -I /home/dtsnko/go/pkg/mod/github.com/grpc-ecosystem/grpc-gateway@v1.16.0/third_party/googleapis  \
      -I /home/dtsnko/go/pkg/mod/github.com/grpc-ecosystem/grpc-gateway@v1.16.0  \
      --go-grpc_out=./proto --go-grpc_opt=paths=source_relative \
      --go_opt=paths=source_relative \
      --go_out=./proto ./proto/*.proto


logger_swagger:
	protoc -I/usr/local/include -I./proto -I /home/dtsnko/go/pkg/mod/github.com/grpc-ecosystem/grpc-gateway@v1.16.0/third_party/googleapis  \
      -I /home/dtsnko/go/pkg/mod/github.com/grpc-ecosystem/grpc-gateway@v1.16.0 \
      --swagger_out=version=false,json_names_for_fields=false,allow_delete_body=true,merge_file_name=logger,allow_merge=true:./swagger \
      ./proto/*.proto