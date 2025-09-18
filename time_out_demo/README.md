# Grpc超时demo

gRPC默认的请求的超时时间是很长的，当你没有设置请求超时时间时，所有在运行的请求都占用大量资源且可能运行很长的时间，导致服务资源损耗过高，使得后来的请求响应过慢，甚至会引起整个进程崩溃。

## 1、编写proto文件 
touch ./proto/simple.proto

## 2、生成go文件
cd proto
protoc --go_out=. --go-grpc_out=. simple.proto  

## 3、编写客户端代码
vi ./client/client.go
### 1.把超时时间设置为当前时间+3秒
```
	clientDeadline := time.Now().Add(time.Duration(3 * time.Second))
	ctx, cancel := context.WithDeadline(ctx, clientDeadline)
	defer cancel()
```
### 2.响应错误检测中添加超时检测
``` 
    // 传入超时时间为3秒的ctx
	res, err := grpcClient.Route(ctx, &req)
	if err != nil {
		//获取错误状态
		statu, ok := status.FromError(err)
		if ok {
			//判断是否为调用超时
			if statu.Code() == codes.DeadlineExceeded {
				log.Fatalln("Route timeout!")
			}
		}
		log.Fatalf("Call Route err: %v", err)
	}
	// 打印返回值
	log.Println(res.Value)
```

## 4、编写服务端代码
vi ./server/server.go
### 当请求超时后，服务端应该停止正在进行的操作，避免资源浪费。
``` 
// Route 实现Route方法
func (s *SimpleService) Route(ctx context.Context, req *pb.SimpleRequest) (*pb.SimpleResponse, error) {
	data := make(chan *pb.SimpleResponse, 1)
	go handle(ctx, req, data)
	select {
	case res := <-data:
		return res, nil
	case <-ctx.Done():
		return nil, status.Errorf(codes.Canceled, "Client cancelled, abandoning.")
	}
}

func handle(ctx context.Context, req *pb.SimpleRequest, data chan<- *pb.SimpleResponse) {
	select {
	case <-ctx.Done():
		log.Println(ctx.Err())
		runtime.Goexit() //超时后退出该Go协程
	case <-time.After(4 * time.Second): // 模拟耗时操作
		res := pb.SimpleResponse{
			Code:  200,
			Value: "hello " + req.Data,
		}
		// //修改数据库前进行超时判断
		// if ctx.Err() == context.Canceled{
		// 	...
		// 	//如果已经超时，则退出
		// }
		data <- &res
	}
}
```
