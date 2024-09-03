package http

import (
	"github.com/skinnykaen/rpa_clone/internal/services"
	"github.com/skinnykaen/rpa_clone/pkg/logger"
)

type Handlers struct {
	ProjectHandler ProjectHandler
	AvatarHandler  AvatarHandler
}

func SetupHandlers(
	loggers logger.Loggers,
	projectService services.ProjectService,
) Handlers {
	return Handlers{
		ProjectHandler: &ProjectHandlerImpl{
			loggers:        loggers,
			projectService: projectService,
		},
		AvatarHandler: &AvatarHandlerImpl{
			loggers: loggers,
		},
	}
}
