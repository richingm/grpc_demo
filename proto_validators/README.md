# proto 验证
## 创建proto文件，添加验证规则
这里使用第三方插件go-proto-validators自动生成验证规则。
```
go get github.com/envoyproxy/protoc-gen-validate
```
1.新建simple.proto文件

2.编译simple.proto文件
```
protoc simple.proto --go_out=. --go-grpc_out=. --validate_out="lang=go:."
```

3、把grpc_validator验证拦截器添加到服务端
```
// 新建gRPC服务器实例
	grpcServer := grpc.NewServer(cred.TLSInterceptor(),
		grpc.StreamInterceptor(grpc_middleware.ChainStreamServer(
			// grpc_ctxtags.StreamServerInterceptor(),
			// grpc_opentracing.StreamServerInterceptor(),
			// grpc_prometheus.StreamServerInterceptor,
			grpc_validator.StreamServerInterceptor(), // proto validate
			grpc_zap.StreamServerInterceptor(zap.ZapInterceptor()),
			grpc_auth.StreamServerInterceptor(auth.AuthInterceptor),
			grpc_recovery.StreamServerInterceptor(recovery.RecoveryInterceptor()),
		)),
		grpc.UnaryInterceptor(grpc_middleware.ChainUnaryServer(
			// grpc_ctxtags.UnaryServerInterceptor(),
			// grpc_opentracing.UnaryServerInterceptor(),
			// grpc_prometheus.UnaryServerInterceptor,
			grpc_validator.UnaryServerInterceptor(), // proto validate
			grpc_zap.UnaryServerInterceptor(zap.ZapInterceptor()),
			grpc_auth.UnaryServerInterceptor(auth.AuthInterceptor),
			grpc_recovery.UnaryServerInterceptor(recovery.RecoveryInterceptor()),
		)),
	)
```
4、客户端请求
```
    // 创建发送结构体
	req := pb.SimpleRequest{
		Name: "test",
		Age:  -1, 
	}
	res, err := grpcClient.GetInfo(context.Background(), &req)
	if err != nil {
		log.Fatalf("Call Route err: %v", err)
	}
```