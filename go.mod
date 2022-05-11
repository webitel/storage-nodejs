module github.com/webitel/storage

go 1.15

require (
	cloud.google.com/go/texttospeech v0.1.0
	github.com/aws/aws-sdk-go v1.43.42
	github.com/go-gorp/gorp v2.2.0+incompatible
	github.com/gorilla/handlers v1.5.1
	github.com/gorilla/mux v1.8.0
	github.com/jmoiron/sqlx v1.3.5
	github.com/lib/pq v1.10.5
	github.com/nicksnyder/go-i18n v1.10.1
	github.com/olivere/elastic v6.2.27+incompatible
	github.com/pborman/uuid v1.2.1
	github.com/pkg/errors v0.9.1
	github.com/robfig/cron v1.2.0
	github.com/webitel/engine v0.0.0-20220511134019-fe681b6fc497
	github.com/webitel/protos/engine v0.0.0-20220511133230-f500804920a4
	github.com/webitel/protos/storage v0.0.0-20220511133230-f500804920a4
	github.com/webitel/wlog v0.0.0-20190823170623-8cc283b29e3e
	google.golang.org/api v0.77.0
	google.golang.org/genproto v0.0.0-20220505152158-f39f71e6c8f3
	google.golang.org/grpc v1.46.0

)

require (
	cloud.google.com/go v0.101.1 // indirect
	github.com/felixge/httpsnoop v1.0.3 // indirect
	github.com/fortytw2/leaktest v1.3.0 // indirect
	github.com/mailru/easyjson v0.7.0 // indirect
	github.com/stretchr/testify v1.7.1 // indirect
)

replace google.golang.org/grpc => google.golang.org/grpc v1.27.0

replace google.golang.org/api => google.golang.org/api v0.54.0
