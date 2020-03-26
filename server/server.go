package server

import (
	"github.com/zjmnssy/etcd"
	"github.com/zjmnssy/serviceRD/health"
	"github.com/zjmnssy/serviceRD/registrar"
	"github.com/zjmnssy/serviceRD/service"
	"github.com/zjmnssy/zlog"
	"google.golang.org/grpc"
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
