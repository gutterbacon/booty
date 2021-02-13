package service

import (
	"github.com/kardianos/service"
	"os/user"
)

type Svc struct {
	service.Service
}

func NewService(config *service.Config) (*Svc, error) {
	u, err := user.Current()
	if err != nil {
		return nil, err
	}
	if u.Uid != "0" {
		// not root, create service as regular user
		config.Option["UserService"] = true
		config.UserName = u.Name
	}
	svc, err := service.New(nil, config)
	return &Svc{svc}, err
}
