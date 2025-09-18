package cred

import (
	"grpc_demo/go-grpc-middleware/pkg/utils"
	"log"
	"os"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

// TLSInterceptor TLS证书认证
func TLSInterceptor() grpc.ServerOption {

	cwd, err := os.Getwd()
	if err != nil {
		log.Fatalf("Failed to get current working directory: %v", err)
	}
	publicFile := utils.GetRealFilePath(cwd, "go-grpc-middleware", "pkg", "tls", "server.pem")
	keyFile := utils.GetRealFilePath(cwd, "go-grpc-middleware", "pkg", "tls", "server.key")
	// 从输入证书文件和密钥文件为服务端构造TLS凭证
	creds, err := credentials.NewServerTLSFromFile(publicFile, keyFile)
	if err != nil {
		log.Fatalf("Failed to generate credentials %v", err)
	}
	return grpc.Creds(creds)
}
