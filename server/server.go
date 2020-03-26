package server

import (
	"github.com/zjmnssy/etcd"
	"github.com/zjmnssy/serviceRD/health"
	"github.com/zjmnssy/serviceRD/registrar"
	"github.com/zjmnssy/serviceRD/service"
	"github.com/zjmnssy/zlog"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

// GetGrpcServer 获取GRPC server
func GetGrpcServer(c etcd.Config, desc service.Desc, serviceName string, ttl int64) (*grpc.Server, *registrar.Registrar, error) {
	impl, err := registrar.NewRegistrar(c, desc, ttl)
	if err != nil {
		zlog.Prints(zlog.Warn, "zgrpc", "create new register error = %s", err)
		return nil, nil, err
	}

	s := grpc.NewServer()

	m := health.GetManager()
	m.Register(s, serviceName)

	return s, impl, nil
}

// GetGrpcServerTLS 获取GRPC TLS server
func GetGrpcServerTLS(c etcd.Config, desc service.Desc, serviceName string, ttl int64, certFile string, keyFile string) (*grpc.Server, *registrar.Registrar, error) {
	impl, err := registrar.NewRegistrar(c, desc, ttl)
	if err != nil {
		zlog.Prints(zlog.Warn, "zgrpc", "create new register error = %s", err)
		return nil, nil, err
	}

	creds, err := credentials.NewServerTLSFromFile(certFile, keyFile)
	if err != nil {
		zlog.Prints(zlog.Warn, "zgrpc", "get creds error = %s", err)
		return nil, nil, err
	}

	s := grpc.NewServer(grpc.Creds(creds))

	m := health.GetManager()
	m.Register(s, serviceName)

	return s, impl, nil
}
