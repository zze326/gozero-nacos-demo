package svc

import "demo/microservice/user/rpc/internal/config"

type ServiceContext struct {
	Config *config.Config
}

func NewServiceContext(c *config.Config) *ServiceContext {
	return &ServiceContext{
		Config: c,
	}
}
