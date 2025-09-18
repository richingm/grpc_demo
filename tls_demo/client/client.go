package main

import (
	"context"
	"grpc_demo/tls_demo/pkg/auth"
	"grpc_demo/tls_demo/pkg/utils"
	"log"
	"os"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"

	pb "grpc_demo/tls_demo/proto"
)

// Address 连接地址
const Address string = ":8000"

var grpcClient pb.SimpleClient

func main() {
	//从输入的证书文件中为客户端构造TLS凭证
	// 获取当前工作目录
	cwd, err := os.Getwd()
	if err != nil {
		log.Fatalf("Failed to get current working directory: %v", err)
	}
	publicFile := utils.GetRealFilePath(cwd, "tls_demo", "pkg", "tls", "server.pem")

	creds, err := credentials.NewClientTLSFromFile(publicFile, "example.com")
	if err != nil {
		log.Fatalf("Failed to create TLS credentials %v", err)
	}
	//构建Token
	token := auth.Token{
		AppID:     "grpc_token",
		AppSecret: "123456",
	}
	// 连接服务器
	conn, err := grpc.NewClient(Address, grpc.WithTransportCredentials(creds), grpc.WithPerRPCCredentials(&token))
	if err != nil {
		log.Fatalf("net.Connect err: %v", err)
	}
	defer func(conn *grpc.ClientConn) {
		err := conn.Close()
		if err != nil {
			panic(err)
		}
	}(conn)

	// 建立gRPC连接
	grpcClient = pb.NewSimpleClient(conn)
	route()
}

// route 调用服务端Route方法
func route() {
	// 创建发送结构体
	req := pb.SimpleRequest{
		Data: "grpc",
	}
	// 调用我们的服务(Route方法)
	// 同时传入了一个 context.Context ，在有需要时可以让我们改变RPC的行为，比如超时/取消一个正在运行的RPC
	res, err := grpcClient.Route(context.Background(), &req)
	if err != nil {
		log.Fatalf("Call Route err: %v", err)
	}
	// 打印返回值
	log.Println(res)
}
