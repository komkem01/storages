package exampletwo

import (
	entitiesinf "storage/app/modules/entities/inf"
	"storage/internal/config"
)

type (
	Module struct {
		Svc *Service
		Ctl *Controller
	}
	Service    struct{}
	Controller struct{}

	Config struct{}
)

func New(conf *config.Config[Config], db entitiesinf.ExampleEntity) *Module {
	return &Module{
		Svc: &Service{},
		Ctl: &Controller{},
	}
}
