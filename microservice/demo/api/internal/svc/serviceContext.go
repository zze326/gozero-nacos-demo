package svc

import (
	commonConfig "demo/common/config"
	"demo/microservice/demo/api/internal/config"
	"demo/microservice/user/rpc/user"
)

type ServiceContext struct {
	Config  *config.Config
	UserRpc user.User
}

func NewServiceContext(c *config.Config, nc *commonConfig.Nacos) *ServiceContext {
	return &ServiceContext{
		Config:  c,
		UserRpc: user.NewUser(nc.NewZrpcClient("user.rpc", c.Name)),
	}
}
