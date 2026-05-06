package storage

import (
	entitiesinf "storage/app/modules/entities/inf"
	"storage/internal/config"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/trace"
)

type Module struct {
	tracer trace.Tracer
	Svc    *Service
	Ctl    *Controller
}

func New(conf *config.Config[Config], ent entitiesinf.StorageEntity) *Module {
	tracer := otel.Tracer("storage.modules.storage")
	svc := newService(&Options{
		Config: conf,
		tracer: tracer,
		ent:    ent,
	})

	return &Module{
		tracer: tracer,
		Svc:    svc,
		Ctl:    newController(tracer, svc),
	}
}
