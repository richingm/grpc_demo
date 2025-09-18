package main

import (
	"context"
	"fmt"
	"google.golang.org/grpc/credentials/insecure"
	pb "grpc_demo/time_out_demo/proto"
	"log"
	"math/rand"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// Address 连接地址
const Address string = ":8000"

var grpcClient pb.SimpleClient

func main() {
	// 连接服务器
	conn, err := grpc.NewClient(Address, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("net.Connect err: %v", err)
	}
	defer func(conn *grpc.ClientConn) {
		err := conn.Close()
		if err != nil {
			panic(err)
		}
	}(conn)

	ctx := context.Background()
	grpcClient = pb.NewSimpleClient(conn)

	// 随机超时
	randomNumber := rand.Intn(3) + 4

	fmt.Println("random number: ", randomNumber)

	route(ctx, time.Duration(randomNumber))
}

// route 调用服务端Route方法
func route(ctx context.Context, deadlines time.Duration) {
	//设置3秒超时时间
	clientDeadline := time.Now().Add(time.Duration(deadlines * time.Second))
	ctx, cancel := context.WithDeadline(ctx, clientDeadline)
	defer cancel()
	// 创建发送结构体
	req := pb.SimpleRequest{
		Data: "grpc",
	}
	// 调用我们的服务(Route方法)
	// 传入超时时间为3秒的ctx
	res, err := grpcClient.Route(ctx, &req)
	if err != nil {
		//获取错误状态
		statusRes, ok := status.FromError(err)
		if ok {
			//判断是否为调用超时
			if statusRes.Code() == codes.DeadlineExceeded {
				log.Fatalln("Route timeout!")
			}
		}
		log.Fatalf("Call Route err: %v", err)
	}
	// 打印返回值
	log.Println(res.Value)
}
