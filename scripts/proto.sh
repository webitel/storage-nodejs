protoc -I/usr/local/include -I../grpc_api/protos -I$GOPATH/src/github.com/grpc-ecosystem/grpc-gateway/third_party/googleapis  \
  -I$GOPATH/src/github.com/grpc-ecosystem/grpc-gateway \
  -I$GOPATH/src/github.com/webitel/protos \
  --swagger_out=version=false,json_names_for_fields=false,allow_delete_body=true,include_package_in_tags=false,allow_repeated_fields_in_body=false,fqn_for_swagger_name=false,merge_file_name=api,allow_merge=true:./api \
  ../grpc_api/protos/*.proto

protoc -I/usr/local/include -I../grpc_api/protos -I$GOPATH/src/github.com/grpc-ecosystem/grpc-gateway/third_party/googleapis  \
  -I/home/igor/opt/Programs/golib/src/github.com/grpc-ecosystem/grpc-gateway \
  -I$GOPATH/src/github.com/webitel/protos \
   --go_out=plugins=grpc,paths=source_relative:../grpc_api/storage ../grpc_api/protos/*.proto
