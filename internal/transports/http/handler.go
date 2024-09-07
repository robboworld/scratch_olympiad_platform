package http

import (
	"github.com/skinnykaen/rpa_clone/internal/services"
	"github.com/skinnykaen/rpa_clone/pkg/logger"
)

type Handlers struct {
	ProjectHandler ProjectHandler
	AvatarHandler  AvatarHandler
	AuthHandler    AuthHandler
}

func SetupHandlers(
	loggers logger.Loggers,
	projectService services.ProjectService,
	authService services.AuthService,
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
	}
}
