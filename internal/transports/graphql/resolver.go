package resolvers

import (
	"github.com/robboworld/scratch_olympiad_platform/internal/services"
	"github.com/robboworld/scratch_olympiad_platform/pkg/logger"
)

type Resolver struct {
	loggers            logger.Loggers
	userService        services.UserService
	authService        services.AuthService
	projectPageService services.ProjectPageService
	settingsService    services.SettingsService
	applicationService services.ApplicationService
	nominationService  services.NominationService
}

func SetupResolvers(
	loggers logger.Loggers,
	userService services.UserService,
	authService services.AuthService,
	projectPageService services.ProjectPageService,
	settingsService services.SettingsService,
	applicationService services.ApplicationService,
	nominationService services.NominationService,
) Resolver {
	return Resolver{
		loggers:            loggers,
		userService:        userService,
		authService:        authService,
		projectPageService: projectPageService,
		settingsService:    settingsService,
		applicationService: applicationService,
		nominationService:  nominationService,
	}
}
