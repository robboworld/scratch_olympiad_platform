package http

import (
	"github.com/robboworld/scratch_olympiad_platform/internal/services"
	"github.com/robboworld/scratch_olympiad_platform/pkg/logger"
)

type Handlers struct {
	ProjectHandler     ProjectHandler
	AvatarHandler      AvatarHandler
	AuthHandler        AuthHandler
	ApplicationHandler ApplicationHandler
	SolutionHandler    SolutionHandler
}

func SetupHandlers(
	loggers logger.Loggers,
	projectService services.ProjectService,
	authService services.AuthService,
	applicationService services.ApplicationService,
) Handlers {
	return Handlers{
		ProjectHandler: ProjectHandler{
			loggers:        loggers,
			projectService: projectService,
		},
		AvatarHandler: AvatarHandler{
			loggers: loggers,
		},
		AuthHandler: AuthHandler{
			loggers:     loggers,
			authService: authService,
		},
		ApplicationHandler: ApplicationHandler{
			loggers:            loggers,
			applicationService: applicationService,
		},
		SolutionHandler: SolutionHandler{
			loggers: loggers,
		},
	}
}
