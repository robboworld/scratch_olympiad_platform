package services

import (
	"github.com/robboworld/scratch_olympiad_platform/internal/api"
	"github.com/robboworld/scratch_olympiad_platform/internal/gateways"
	"go.uber.org/fx"
)

type Services struct {
	fx.Out
	UserService        UserService
	AuthService        AuthService
	ProjectService     ProjectService
	ProjectPageService ProjectPageService
	SettingsService    SettingsService
	ApplicationService ApplicationService
	NominationService  NominationService
	CountryService     CountryService
}

func SetupServices(
	userGateway gateways.UserGateway,
	projectGateway gateways.ProjectGateway,
	projectPageGateway gateways.ProjectPageGateway,
	settingsGateway gateways.SettingsGateway,
	applicationGateway gateways.ApplicationGateway,
	nominationGateway gateways.NominationGateway,
	countryGateway gateways.CountryGateway,
	applicationAPI api.ApplicationAPI,
) Services {
	return Services{
		UserService: &UserServiceImpl{
			userGateway:    userGateway,
			countryGateway: countryGateway,
		},
		AuthService: &AuthServiceImpl{
			userGateway:     userGateway,
			countryGateway:  countryGateway,
			settingsGateway: settingsGateway,
		},
		ProjectService: &ProjectServiceImpl{
			projectGateway: projectGateway,
		},
		ProjectPageService: &ProjectPageServiceImpl{
			projectGateway:     projectGateway,
			projectPageGateway: projectPageGateway,
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
	}
}
