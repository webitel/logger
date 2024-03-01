.PHONY: all

all: update_proto generate_proto generate_swagger


update_proto:
	protodep up --use-https

generate_proto:
	protoc -I/usr/local/include -I./proto -I./  \
			-I./proto -I /home/dtsnko/go/pkg/mod/github.com/grpc-ecosystem/grpc-gateway@v1.16.0/third_party/googleapis  \
          --go-grpc_out=./api --go-grpc_opt=paths=source_relative \
          --go_opt=paths=source_relative \
          --go_out=./api ./proto/*.proto


generate_swagger:
	protoc -I/usr/local/include -I./proto -I /home/dtsnko/go/pkg/mod/github.com/grpc-ecosystem/grpc-gateway@v1.16.0/third_party/googleapis  \
      -I /home/dtsnko/go/pkg/mod/github.com/grpc-ecosystem/grpc-gateway@v1.16.0 \
      --swagger_out=version=false,json_names_for_fields=false,allow_delete_body=true,merge_file_name=logger,allow_merge=true:./swagger \
      ./proto/*.proto