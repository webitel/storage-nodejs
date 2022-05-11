export SPEECHSDK_ROOT="$HOME/speechsdk"
export CGO_CFLAGS="-I$HOME/speechsdk/include/c_api"
export CGO_LDFLAGS="-L$HOME/speechsdk/lib/x64 -lMicrosoft.CognitiveServices.Speech.core"
export LD_LIBRARY_PATH="$HOME/speechsdk/lib/x64:$LD_LIBRARY_PATH"
go run app.go -id='igor' \
-data_source='postgres://opensips:webitel@10.9.8.111:5432/webitel?fallback_application_name=engine&sslmode=disable&connect_timeout=10&search_path=call_center' \
-consul='10.9.8.111:8500' \
-file_store_type \
"local" \
-file_store_props \
"{\"directory\":\"/tmp\",\"path_pattern\": \"$DOMAIN/$Y/$M/$D/$H\"}" \
-presigned_cert \
"/home/igor/work/storage/bin/key.pem"