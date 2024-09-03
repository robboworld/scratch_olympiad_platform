package server

import (
	"context"
	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/rs/cors"
	"github.com/skinnykaen/rpa_clone/graph"
	"github.com/skinnykaen/rpa_clone/internal/consts"
	"github.com/skinnykaen/rpa_clone/internal/graphql/directives"
	resolvers "github.com/skinnykaen/rpa_clone/internal/transports/graphql"
	http2 "github.com/skinnykaen/rpa_clone/internal/transports/http"
	"github.com/skinnykaen/rpa_clone/pkg/logger"
	"github.com/spf13/viper"
	"go.uber.org/fx"
	"net/http"
)

func NewServer(
	m consts.Mode,
	lifecycle fx.Lifecycle,
	loggers logger.Loggers,
	resolver resolvers.Resolver,
	handlers http2.Handlers,
) {
	lifecycle.Append(
		fx.Hook{
			OnStart: func(ctx context.Context) (err error) {
				port := viper.GetString("graphql_server_port")
				c := graph.Config{Resolvers: &resolver}
				c.Directives.HasRole = directives.HasRole(loggers.Err)
				mux := http.NewServeMux()
				srv := handler.NewDefaultServer(graph.NewExecutableSchema(c))
				switch m {
				case consts.Production:
					mux.Handle("/query", Auth(srv, loggers.Err))
					mux.Handle("/project", Auth(handlers.ProjectHandler, loggers.Err))
				case consts.Development:
					mux.Handle("/", playground.Handler("GraphQL playground", "/query"))
					mux.Handle("/query", Auth(srv, loggers.Err))
					mux.Handle("/project", Auth(handlers.ProjectHandler, loggers.Err))
					mux.Handle("/avatar", Auth(handlers.AvatarHandler, loggers.Err))
				}
				loggers.Info.Printf(
					"Connect to %s:%s/ for GraphQL playground",
					viper.GetString("server_host"),
					port,
				)
				loggers.Info.Printf(
					"The app is running in %s mode",
					m,
				)
				go func() {
					loggers.Err.Fatal(http.ListenAndServe(":"+port, cors.New(
						cors.Options{
							AllowedOrigins:   viper.GetStringSlice("cors.allowed_origins"),
							AllowCredentials: viper.GetBool("cors.allow_credentials"),
							AllowedMethods:   viper.GetStringSlice("cors.allowed_methods"),
							AllowedHeaders:   viper.GetStringSlice("cors.allowed_headers"),
						},
					).Handler(mux)))
				}()
				return
			},
			OnStop: func(context.Context) error {
				return nil
			},
		})
}
