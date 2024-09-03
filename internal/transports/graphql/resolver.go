package resolvers

import (
	"github.com/skinnykaen/rpa_clone/internal/services"
	"github.com/skinnykaen/rpa_clone/pkg/logger"
)

type Resolver struct {
	loggers            logger.Loggers
	userService        services.UserService
	authService        services.AuthService
	projectPageService services.ProjectPageService
	settingsService    services.SettingsService
}

func SetupResolvers(
	loggers logger.Loggers,
	userService services.UserService,
	authService services.AuthService,
	projectPageService services.ProjectPageService,
	settingsService services.SettingsService,
) Resolver {
	return Resolver{
		loggers:            loggers,
		userService:        userService,
		authService:        authService,
		projectPageService: projectPageService,
		settingsService:    settingsService,
	}
}
