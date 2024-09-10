package gateways

import (
	"github.com/robboworld/scratch_olympiad_platform/internal/db"
	"go.uber.org/fx"
)

type Gateways struct {
	fx.Out
	UserGateway UserGateway
	ParentRel   ParentRel
	Project     ProjectGateway
	ProjectPage ProjectPageGateway
	Settings    SettingsGateway
	Application ApplicationGateway
	Nomination  NominationGateway
	Country     CountryGateway
}

func SetupGateways(pc db.PostgresClient) Gateways {
	return Gateways{
		UserGateway: UserGatewayImpl{pc},
		ParentRel:   ParentRelGatewayImpl{pc},
		Project:     ProjectGatewayImpl{pc},
		ProjectPage: ProjectPageGatewayImpl{pc},
		Settings:    SettingsGatewayImpl{pc},
		Application: ApplicationGatewayImpl{pc},
		Nomination:  NominationGatewayImpl{pc},
		Country:     CountryGatewayImpl{pc},
	}
}
