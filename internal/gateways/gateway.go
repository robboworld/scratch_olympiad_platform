package gateways

import (
	"github.com/skinnykaen/rpa_clone/internal/db"
	"go.uber.org/fx"
)

type Gateways struct {
	fx.Out
	UserGateway UserGateway
	ParentRel   ParentRel
	Project     ProjectGateway
	ProjectPage ProjectPageGateway
	Settings    SettingsGateway
}

func SetupGateways(pc db.PostgresClient) Gateways {
	return Gateways{
		UserGateway: UserGatewayImpl{pc},
		ParentRel:   ParentRelGatewayImpl{pc},
		Project:     ProjectGatewayImpl{pc},
		ProjectPage: ProjectPageGatewayImpl{pc},
		Settings:    SettingsGatewayImpl{pc},
	}
}
