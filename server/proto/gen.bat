protoc -I =. --go_out=plugins=grpc,paths=source_relative:gen/go trip.proto
@REM grpc-gateway是工具名 _out是对工具参数的配置
protoc -I =. --grpc-gateway_out=paths=source_relative,grpc_api_configuration=trip.yaml:gen/go trip.proto