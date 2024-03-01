.PHONY: all

all: update_native_proto update_storage_proto native_go storage_go generate_swagger



update_native_proto:
	$(MAKE) -C proto/native all

update_storage_proto:
	$(MAKE) -C proto/storage all



native_go:
	protoc  -I./proto/native/contracts  \
			-I /home/dtsnko/go/pkg/mod/github.com/grpc-ecosystem/grpc-gateway@v1.16.0/third_party/googleapis  \
			--proto_path=./proto/native/contracts \
          --go-grpc_out=./api/native --go-grpc_opt=paths=source_relative \
          --go_opt=paths=source_relative \
          --go_out=./api/native ./proto/native/contracts/*.proto

  storage_go:
	protoc  -I /home/dtsnko/go/pkg/mod/github.com/grpc-ecosystem/grpc-gateway@v1.16.0/third_party/googleapis  \
			--go-grpc_out=./api/storage --go-grpc_opt=paths=source_relative \
			--proto_path=./proto/storage/contracts --go_out=./api/storage --go_opt=paths=source_relative ./proto/storage/contracts/file.proto


generate_swagger:
	protoc  -I /home/dtsnko/go/pkg/mod/github.com/grpc-ecosystem/grpc-gateway@v1.16.0/third_party/googleapis  \
      -I /home/dtsnko/go/pkg/mod/github.com/grpc-ecosystem/grpc-gateway@v1.16.0 \
      --swagger_out=version=false,json_names_for_fields=false,allow_delete_body=true,merge_file_name=logger,allow_merge=true:./api/native \
      --proto_path=./proto/native/contracts \
      ./proto/native/contracts/*.proto