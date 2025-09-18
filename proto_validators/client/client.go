package main

import (
	"context"
	"grpc_demo/proto_validators/client/auth"
	"grpc_demo/proto_validators/pkg/utils"
	"log"
	"os"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"

	pb "grpc_demo/proto_validators/proto"
)

// Address 连接地址
const Address string = ":8000"

func main() {
	cwd, err := os.Getwd()
	if err != nil {
		log.Fatalf("Failed to get current working directory: %v", err)
	}
	publicFile := utils.GetRealFilePath(cwd, "proto_validators", "pkg", "tls", "server.pem")

	//从输入的证书文件中为客户端构造TLS凭证
	creds, err := credentials.NewClientTLSFromFile(publicFile, "example.com")
	if err != nil {
		log.Fatalf("Failed to create TLS credentials %v", err)
	}
	//构建Token
	token := auth.Token{
		Value: "bearer grpc.auth.token",
	}
	// 连接服务器
	conn, err := grpc.NewClient(Address, grpc.WithTransportCredentials(creds), grpc.WithPerRPCCredentials(&token))
	if err != nil {
		log.Fatalf("net.Connect err: %v", err)
	}
	defer conn.Close()

	// 建立gRPC连接
	grpcClient := pb.NewSimpleServiceClient(conn)

	// 创建发送结构体
	req := pb.SimpleRequest{
		Name: "test",
		Age:  -1,
	}
	res, err := grpcClient.GetInfo(context.Background(), &req)
	if err != nil {
		log.Fatalf("Call Route err: %v", err)
	}
	// 打印返回值
	log.Println(res)
}
