package collector

import "storage/internal/config"

type Module struct {
	Svc *Service
}

func New(conf *config.Config[Config]) *Module {
	return &Module{
		Svc: newService(conf),
	}
}
