package services

import (
	"github.com/robboworld/scratch_olympiad_platform/internal/api"
	"github.com/robboworld/scratch_olympiad_platform/internal/gateways"
	"go.uber.org/fx"
)

type Services struct {
	fx.Out
	UserService             UserService
	AuthService             AuthService
	SettingsService         SettingsService
	ApplicationService      ApplicationService
	NominationService       NominationService
	CountryService          CountryService
	RegionService           RegionService
	SolutionService         SolutionService
	EventService            EventService
	EventTranslationService EventTranslationService
	EventCountryService     EventCountryService
	EventRegionService      EventRegionService
	EventUserService        EventUserService
}

func SetupServices(
	userGateway gateways.UserGateway,
	authDataGateway gateways.AuthDataGateway,
	settingsGateway gateways.SettingsGateway,
	applicationGateway gateways.ApplicationGateway,
	nominationGateway gateways.NominationGateway,
	countryGateway gateways.CountryGateway,
	regionGateway gateways.RegionGateway,
	solutionGateway gateways.SolutionGateway,
	eventGateway gateways.EventGateway,
	eventTranslationGateway gateways.EventTranslationGateway,
	eventCountryGateway gateways.EventCountryRelGateway,
	eventRegionGateway gateways.EventRegionRelGateway,
	eventUserGateway gateways.EventUserRelGateway,
	applicationAPI api.ApplicationAPI,
) Services {
	return Services{
		UserService: &UserServiceImpl{
			userGateway:    userGateway,
			countryGateway: countryGateway,
			regionGateway:  regionGateway,
		},
		AuthService: &AuthServiceImpl{
			userGateway:     userGateway,
			authDataGateway: authDataGateway,
			countryGateway:  countryGateway,
			regionGateway:   regionGateway,
			settingsGateway: settingsGateway,
		},
		SettingsService: &SettingsServiceImpl{
			settingsGateway: settingsGateway,
		},
		ApplicationService: &ApplicationServiceImpl{
			applicationGateway: applicationGateway,
			nominationGateway:  nominationGateway,
			userGateway:        userGateway,
			applicationAPI:     applicationAPI,
		},
		NominationService: &NominationServiceImpl{
			nominationGateway: nominationGateway,
		},
		CountryService: &CountryServiceImpl{
			countryGateway: countryGateway,
		},
		RegionService: &RegionServiceImpl{
			regionGateway: regionGateway,
		},
		SolutionService: &SolutionServiceImpl{
			solutionGateway: solutionGateway,
			userGateway:     userGateway,
		},
		EventService: &EventServiceImpl{
			userGateway:             userGateway,
			eventGateway:            eventGateway,
			eventUserGateway:        eventUserGateway,
			eventTranslationGateway: eventTranslationGateway,
			eventCountryGateway:     eventCountryGateway,
			eventRegionGateway:      eventRegionGateway,
		},
		EventTranslationService: &EventTranslationServiceImpl{
			eventGateway:            eventGateway,
			eventTranslationGateway: eventTranslationGateway,
			eventUserGateway:        eventUserGateway,
		},
		EventCountryService: &EventCountryServiceImpl{
			countryGateway:      countryGateway,
			eventGateway:        eventGateway,
			eventCountryGateway: eventCountryGateway,
		},
		EventRegionService: &EventRegionServiceImpl{
			regionGateway:       regionGateway,
			eventGateway:        eventGateway,
			eventCountryGateway: eventCountryGateway,
			eventRegionGateway:  eventRegionGateway,
		},
		EventUserService: &EventUserServiceImpl{
			userGateway:         userGateway,
			eventGateway:        eventGateway,
			eventUserGateway:    eventUserGateway,
			eventCountryGateway: eventCountryGateway,
			eventRegionGateway:  eventRegionGateway,
		},
	}
}
