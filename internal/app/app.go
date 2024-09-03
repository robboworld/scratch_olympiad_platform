package app

import (
	"github.com/skinnykaen/rpa_clone/internal/configs"
	"github.com/skinnykaen/rpa_clone/internal/consts"
	"github.com/skinnykaen/rpa_clone/internal/db"
	"github.com/skinnykaen/rpa_clone/internal/gateways"
	"github.com/skinnykaen/rpa_clone/internal/server"
	"github.com/skinnykaen/rpa_clone/internal/services"
	resolvers "github.com/skinnykaen/rpa_clone/internal/transports/graphql"
	"github.com/skinnykaen/rpa_clone/internal/transports/http"
	"github.com/skinnykaen/rpa_clone/pkg/logger"
	"go.uber.org/fx"
	"log"
	"os"
)

func InvokeWith(m consts.Mode, options ...fx.Option) *fx.App {
	if err := configs.Init(m); err != nil {
		log.Fatalf("%s", err.Error())
	}
	di := []fx.Option{
		fx.Provide(func() consts.Mode { return m }),
		fx.Provide(logger.InitLogger),
		fx.Provide(db.InitPostgresClient),
		fx.Provide(gateways.SetupGateways),
		fx.Provide(services.SetupServices),
		fx.Provide(resolvers.SetupResolvers),
		fx.Provide(http.SetupHandlers),
	}
	for _, option := range options {
		di = append(di, option)
	}
	return fx.New(di...)
}

func RunApp() {
	if len(os.Args) == 2 && (consts.Mode(os.Args[1]) == consts.Development ||
		consts.Mode(os.Args[1]) == consts.Production) {
		InvokeWith(consts.Mode(os.Args[1]), fx.Invoke(server.NewServer)).Run()
	} else {
		InvokeWith(consts.Development, fx.Invoke(server.NewServer)).Run()
	}
}
