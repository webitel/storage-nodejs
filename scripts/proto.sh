protoc -I/usr/local/include -I../grpc_api/proto -I$GOPATH/src/github.com/grpc-ecosystem/grpc-gateway/third_party/googleapis  \
  -I/home/igor/opt/Programs/golib/src/github.com/grpc-ecosystem/grpc-gateway \
  -I$GOPATH/src/github.com/webitel/engine/grpc_api/protos \
  --swagger_out=version=false,json_names_for_fields=false,allow_delete_body=true,include_package_in_tags=false,allow_repeated_fields_in_body=false,fqn_for_swagger_name=false,merge_file_name=api,allow_merge=true:./api \
  ../grpc_api/proto/*.proto

protoc -I/usr/local/include -I../grpc_api/proto -I$GOPATH/src/github.com/grpc-ecosystem/grpc-gateway/third_party/googleapis  \
  -I/home/igor/opt/Programs/golib/src/github.com/grpc-ecosystem/grpc-gateway \
  -I$GOPATH/src/github.com/webitel/engine/grpc_api/protos \
   --go_out=plugins=grpc:../grpc_api/storage ../grpc_api/proto/*.proto
