package app

import (
	"github.com/robboworld/scratch_olympiad_platform/internal/api"
	"github.com/robboworld/scratch_olympiad_platform/internal/configs"
	"github.com/robboworld/scratch_olympiad_platform/internal/consts"
	"github.com/robboworld/scratch_olympiad_platform/internal/db"
	"github.com/robboworld/scratch_olympiad_platform/internal/gateways"
	"github.com/robboworld/scratch_olympiad_platform/internal/server"
	"github.com/robboworld/scratch_olympiad_platform/internal/services"
	resolvers "github.com/robboworld/scratch_olympiad_platform/internal/transports/graphql"
	"github.com/robboworld/scratch_olympiad_platform/internal/transports/http"
	"github.com/robboworld/scratch_olympiad_platform/pkg/logger"
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
		fx.Provide(api.SetupAPI),
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
