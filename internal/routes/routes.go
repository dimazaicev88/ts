package routes

import (
	"context"

	"github.com/dimazaicev88/ts/config"
	httpHandlers "github.com/dimazaicev88/ts/internal/handlers/http"
	"github.com/dimazaicev88/ts/internal/services"

	"github.com/gofiber/contrib/v3/swaggo"
	"github.com/gofiber/fiber/v3"
)

type Routes struct {
	metadata httpHandlers.Metadata
	tasks    httpHandlers.Tasks
	workers  httpHandlers.Workers
	fb       *fiber.App
}

func NewRoutes(
	ctx context.Context,
	fb *fiber.App,
	cfg config.ServerConfig,
	allServices services.AllServices,
) *Routes {
	return &Routes{
		fb:       fb,
		metadata: httpHandlers.NewMetadata(ctx, fb, cfg),
		tasks:    httpHandlers.NewTasks(ctx, fb, allServices.TaskService),
		workers:  httpHandlers.NewWorkers(ctx, fb, allServices.WorkerService),
	}
}

func (p Routes) AddRoutes() {
	p.fb.Get("/api/v1/server/metadata", p.metadata.GetMetadata)
	p.fb.Get("/swagger/*", swaggo.HandlerDefault)

	//Task API
	p.fb.Get("/api/v1/tasks/workers/:workerUid", p.tasks.FindTask)
	p.fb.Delete("/api/v1/tasks", p.tasks.DeleteTask)
	p.fb.Post("/api/v1/tasks/stop", p.tasks.StopTask)

	//Worker API
	p.fb.Get("/api/v1/wrokers/:workerName", p.workers.FindByName)
	p.fb.Get("/api/v1/wrokers", p.workers.Find)
}
