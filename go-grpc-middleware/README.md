# 中间件

上篇介绍了gRPC中TLS认证和自定义方法认证，最后还简单介绍了gRPC拦截器的使用。gRPC自身只能设置一个拦截器，所有逻辑都写一起会比较乱

go-grpc-middleware封装了认证（auth）, 日志（ logging）, 消息（message）, 验证（validation）, 重试（retries） 和监控（retries）等拦截器。

安装: go get github.com/grpc-ecosystem/go-grpc-middleware

使用:
``` 
import "github.com/grpc-ecosystem/go-grpc-middleware"
myServer := grpc.NewServer(
    grpc.StreamInterceptor(grpc_middleware.ChainStreamServer(
        grpc_ctxtags.StreamServerInterceptor(),
        grpc_opentracing.StreamServerInterceptor(),
        grpc_prometheus.StreamServerInterceptor,
        grpc_zap.StreamServerInterceptor(zapLogger),
        grpc_auth.StreamServerInterceptor(myAuthFunction),
        grpc_recovery.StreamServerInterceptor(),
    )),
    grpc.UnaryInterceptor(grpc_middleware.ChainUnaryServer(
        grpc_ctxtags.UnaryServerInterceptor(),
        grpc_opentracing.UnaryServerInterceptor(),
        grpc_prometheus.UnaryServerInterceptor,
        grpc_zap.UnaryServerInterceptor(zapLogger),
        grpc_auth.UnaryServerInterceptor(myAuthFunction),
        grpc_recovery.UnaryServerInterceptor(),
    )),
)
# grpc.StreamInterceptor中添加流式RPC的拦截器。  
# grpc.UnaryInterceptor中添加简单RPC的拦截器。 
```

### grpc_zap日志记录
1.创建zap.Logger实例
```
func ZapInterceptor() *zap.Logger {
	logger, err := zap.NewDevelopment()
	if err != nil {
		log.Fatalf("failed to initialize zap logger: %v", err)
	}
	grpc_zap.ReplaceGrpcLogger(logger)
	return logger
}
```
2.把zap拦截器添加到服务端
``` 
grpcServer := grpc.NewServer(
	grpc.StreamInterceptor(grpc_middleware.ChainStreamServer(
			grpc_zap.StreamServerInterceptor(zap.ZapInterceptor()),
		)),
		grpc.UnaryInterceptor(grpc_middleware.ChainUnaryServer(
			grpc_zap.UnaryServerInterceptor(zap.ZapInterceptor()),
		)),
	)
```
3.查看日志 log/debug.log
4.把日志写到文件中
现在我们把日志写到文件中，修改ZapInterceptor方法。
```
import (
	grpc_zap "github.com/grpc-ecosystem/go-grpc-middleware/logging/zap"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

// ZapInterceptor 返回zap.logger实例(把日志写到文件中)
func ZapInterceptor() *zap.Logger {
	w := zapcore.AddSync(&lumberjack.Logger{
		Filename:  "log/debug.log",
		MaxSize:   1024, //MB
		LocalTime: true,
	})

	config := zap.NewProductionEncoderConfig()
	config.EncodeTime = zapcore.ISO8601TimeEncoder
	core := zapcore.NewCore(
		zapcore.NewJSONEncoder(config),
		w,
		zap.NewAtomicLevel(),
	)

	logger := zap.New(core, zap.AddCaller(), zap.AddCallerSkip(1))
	grpc_zap.ReplaceGrpcLogger(logger)
	return logger
}
```






