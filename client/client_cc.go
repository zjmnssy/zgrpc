package client

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/zjmnssy/serviceRD/balancer"
	"github.com/zjmnssy/zlog"
	"google.golang.org/grpc"

	"google.golang.org/grpc/credentials"

	// 加载健康检查库
	_ "google.golang.org/grpc/health"
)

// NameUnit 服务方法
type NameUnit struct {
	Service string `json:"service"` // package.Service
	//Method  string `json:"method"`
}

// RetryPolicy 重试策略
type RetryPolicy struct {
	MaxAttempts          int    `json:"maxAttempts"` // >= 2 包括正常调用的第一次，所以如果等于２就表示重试一次
	InitialBackoff       string `json:"initialBackoff"`
	MaxBackoff           string `json:"maxBackoff"`
	BackoffMultiplier    int    `json:"backoffMultiplier"`    // > 0
	RetryableStatusCodes []int  `json:"retryableStatusCodes"` // eg. [14]
}

// MethodConfigUnit 方法配置
type MethodConfigUnit struct {
	Name                    []NameUnit  `json:"name"`
	RetryPolicy             RetryPolicy `json:"retryPolicy"`
	WaitForReady            bool        `json:"waitForReady"`
	Timeout                 string      `json:"timeout"`
	MaxRequestMessageBytes  int         `json:"maxRequestMessageBytes"`
	MaxResponseMessageBytes int         `json:"maxResponseMessageBytes"`
}

// RetryThrottling 重试阈值控制
type RetryThrottling struct {
	MaxTokens  uint `json:"maxTokens"`  // (0, 1000]
	TokenRatio uint `json:"tokenRatio"` // (0, 1]
}

// HealthCheckConfig 健康检查配置
type HealthCheckConfig struct {
	ServiceName string `json:"serviceName"` // package.Service
}

// ServerConfig grpc服务配置
type ServerConfig struct {
	LoadBalancingPolicy string             `json:"loadBalancingPolicy"`
	MethodConfig        []MethodConfigUnit `json:"methodConfig"`
	RetryThrottling     RetryThrottling    `json:"retryThrottling"`
	HealthCheckConfig   HealthCheckConfig  `json:"healthCheckConfig"`
}

// GetGrpcConn 获取GRPC client连接
func GetGrpcConn(serviceName string, scheme string) (*grpc.ClientConn, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(2)*time.Second)
	defer cancel()

	var serverConfig ServerConfig
	serverConfig.LoadBalancingPolicy = balancer.Random
	var methodConfig = MethodConfigUnit{Name: make([]NameUnit, 0, 0)}
	var nameUnit = NameUnit{Service: serviceName}
	methodConfig.Name = append(methodConfig.Name, nameUnit)
	var retryPolicy = RetryPolicy{RetryableStatusCodes: make([]int, 0, 0)}
	retryPolicy.RetryableStatusCodes = append(retryPolicy.RetryableStatusCodes, 14)
	retryPolicy.MaxAttempts = 3
	retryPolicy.InitialBackoff = "0.1s"
	retryPolicy.MaxBackoff = "1s"
	retryPolicy.BackoffMultiplier = 1
	methodConfig.RetryPolicy = retryPolicy
	methodConfig.WaitForReady = true
	methodConfig.Timeout = "1.5s"
	methodConfig.MaxRequestMessageBytes = 1024 * 1024 * 1024
	methodConfig.MaxResponseMessageBytes = 1024 * 1024 * 1024
	serverConfig.MethodConfig = make([]MethodConfigUnit, 0, 0)
	serverConfig.MethodConfig = append(serverConfig.MethodConfig, methodConfig)
	var retryThrottling = RetryThrottling{MaxTokens: 1000, TokenRatio: 1}
	serverConfig.RetryThrottling = retryThrottling
	var healthCheckConfig = HealthCheckConfig{ServiceName: serviceName}
	serverConfig.HealthCheckConfig = healthCheckConfig

	bytes, err := json.Marshal(serverConfig)
	if err != nil {
		return nil, err
	}

	cc, err := grpc.DialContext(ctx,
		fmt.Sprintf("%s:///", scheme),
		//grpc.WithBlock(), // 如果使用WithBlock()， 此接口返回失败，导致外面调用不好处理， 可能进入不了服务发现和负载均衡
		grpc.WithInsecure(),
		grpc.WithBackoffMaxDelay(time.Second),
		grpc.WithDisableServiceConfig(),
		grpc.WithDefaultServiceConfig(string(bytes)),
	)
	if err != nil {
		zlog.Prints(zlog.Warn, "zgrpc", "grpc dial error = %s", err)
		return nil, err
	}

	return cc, nil
}

// GetGrpcConnTLS 获取GRPC TLS client连接
func GetGrpcConnTLS(serviceName string, scheme string, certFile string, keyFile string) (*grpc.ClientConn, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(2)*time.Second)
	defer cancel()

	var serverConfig ServerConfig
	serverConfig.LoadBalancingPolicy = balancer.Random
	var methodConfig = MethodConfigUnit{Name: make([]NameUnit, 0, 0)}
	var nameUnit = NameUnit{Service: serviceName}
	methodConfig.Name = append(methodConfig.Name, nameUnit)
	var retryPolicy = RetryPolicy{RetryableStatusCodes: make([]int, 0, 0)}
	retryPolicy.RetryableStatusCodes = append(retryPolicy.RetryableStatusCodes, 14)
	retryPolicy.MaxAttempts = 3
	retryPolicy.InitialBackoff = "0.1s"
	retryPolicy.MaxBackoff = "1s"
	retryPolicy.BackoffMultiplier = 1
	methodConfig.RetryPolicy = retryPolicy
	methodConfig.WaitForReady = true
	methodConfig.Timeout = "1.5s"
	methodConfig.MaxRequestMessageBytes = 1024 * 1024 * 1024
	methodConfig.MaxResponseMessageBytes = 1024 * 1024 * 1024
	serverConfig.MethodConfig = make([]MethodConfigUnit, 0, 0)
	serverConfig.MethodConfig = append(serverConfig.MethodConfig, methodConfig)
	var retryThrottling = RetryThrottling{MaxTokens: 1000, TokenRatio: 1}
	serverConfig.RetryThrottling = retryThrottling
	var healthCheckConfig = HealthCheckConfig{ServiceName: serviceName}
	serverConfig.HealthCheckConfig = healthCheckConfig

	bytes, err := json.Marshal(serverConfig)
	if err != nil {
		return nil, err
	}

	creds, err := credentials.NewServerTLSFromFile(certFile, keyFile)
	if err != nil {
		zlog.Prints(zlog.Warn, "zgrpc", "get creds error = %s", err)
		return nil, err
	}

	cc, err := grpc.DialContext(ctx,
		fmt.Sprintf("%s:///", scheme),
		//grpc.WithBlock(), // 如果使用WithBlock()， 此接口返回失败，导致外面调用不好处理， 可能进入不了服务发现和负载均衡
		grpc.WithInsecure(),
		grpc.WithBackoffMaxDelay(time.Second),
		grpc.WithDisableServiceConfig(),
		grpc.WithTransportCredentials(creds),
		grpc.WithDefaultServiceConfig(string(bytes)),
	)
	if err != nil {
		zlog.Prints(zlog.Warn, "zgrpc", "grpc dial error = %s", err)
		return nil, err
	}

	return cc, nil
}
