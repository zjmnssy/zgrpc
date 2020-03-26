package server

import (
	"context"
	"fmt"
	"net"

	"google.golang.org/grpc/peer"
)

// GetClientIP 获取请求客户端的远程地址, 通过从metadata中获取远程地址信息
func GetClientIP(ctx context.Context) (string, error) {
	pr, ok := peer.FromContext(ctx)
	if !ok {
		return "", fmt.Errorf("invoke FromContext() failed")
	}

	if pr.Addr == net.Addr(nil) {
		return "", fmt.Errorf("peer.Addr is nil")
	}

	return pr.Addr.String(), nil
}
