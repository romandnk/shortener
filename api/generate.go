package api

//go:generate protoc --go_out=../internal/server/grpc/url/pb --go-grpc_out=../internal/server/grpc/url/pb url/URLService.proto
