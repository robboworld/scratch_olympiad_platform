package resolvers

import (
	"github.com/robboworld/scratch_olympiad_platform/internal/services"
	"github.com/robboworld/scratch_olympiad_platform/pkg/logger"
)

type Resolver struct {
	loggers                 logger.Loggers
	userService             services.UserService
	authService             services.AuthService
	settingsService         services.SettingsService
	applicationService      services.ApplicationService
	nominationService       services.NominationService
	countryService          services.CountryService
	regionService           services.RegionService
	eventService            services.EventService
	eventTranslationService services.EventTranslationService
	eventCountryService     services.EventCountryService
	eventRegionService      services.EventRegionService
	eventUserService        services.EventUserService
}

func SetupResolvers(
	loggers logger.Loggers,
	userService services.UserService,
	authService services.AuthService,
	settingsService services.SettingsService,
	applicationService services.ApplicationService,
	nominationService services.NominationService,
	countryService services.CountryService,
	regionService services.RegionService,
	eventService services.EventService,
	eventTranslationService services.EventTranslationService,
	eventCountryService services.EventCountryService,
	eventRegionService services.EventRegionService,
	eventUserService services.EventUserService,
) Resolver {
	return Resolver{
		loggers:                 loggers,
		userService:             userService,
		authService:             authService,
		settingsService:         settingsService,
		applicationService:      applicationService,
		nominationService:       nominationService,
		countryService:          countryService,
		regionService:           regionService,
		eventService:            eventService,
		eventTranslationService: eventTranslationService,
		eventCountryService:     eventCountryService,
		eventRegionService:      eventRegionService,
		eventUserService:        eventUserService,
	}
}
