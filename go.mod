module github.com/richingm/grpc_demo

go 1.23.6

require (
	github.com/grpc-ecosystem/go-grpc-middleware v1.4.0
	go.uber.org/zap v1.27.0
	google.golang.org/grpc v1.75.1
	google.golang.org/protobuf v1.36.9
	gopkg.in/natefinch/lumberjack.v2 v2.2.1
)

require (
	cloud.google.com/go/compute/metadata v0.7.0 // indirect
	github.com/golang/protobuf v1.5.4 // indirect
	go.uber.org/multierr v1.10.0 // indirect
	golang.org/x/net v0.41.0 // indirect
	golang.org/x/sys v0.33.0 // indirect
	golang.org/x/text v0.26.0 // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20250707201910-8d1bb00bc6a7 // indirect
)

replace cloud.google.com/go/compute/metadata => cloud.google.com/go/compute/metadata v0.7.0
