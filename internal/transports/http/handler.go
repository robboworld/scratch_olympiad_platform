package http

import (
	"github.com/robboworld/scratch_olympiad_platform/internal/services"
	"github.com/robboworld/scratch_olympiad_platform/pkg/logger"
)

type Handlers struct {
	AvatarHandler   AvatarHandler
	SolutionHandler SolutionHandler
}

func SetupHandlers(
	loggers logger.Loggers,
	solutionService services.SolutionService,
) Handlers {
	return Handlers{
		AvatarHandler: AvatarHandler{
			loggers: loggers,
		},
		SolutionHandler: SolutionHandler{
			loggers:         loggers,
			solutionService: solutionService,
		},
	}
}
