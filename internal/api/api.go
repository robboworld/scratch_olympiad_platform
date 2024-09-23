package api

import (
	"go.uber.org/fx"
)

type API struct {
	fx.Out
	ApplicationAPI ApplicationAPI
}

func SetupAPI() API {
	return API{
		ApplicationAPI: NewApplicationAPIImpl(),
	}
}
