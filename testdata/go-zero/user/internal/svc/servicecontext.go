package svc

import (
	"github.com/CloverOS/goctl-swag/testdata/go-zero/user/internal/config"
)

type ServiceContext struct {
	Config config.Config
}

func NewServiceContext(c config.Config) *ServiceContext {
	return &ServiceContext{
		Config: c,
	}
}
