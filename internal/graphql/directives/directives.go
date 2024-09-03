package directives

import (
	"context"
	"github.com/99designs/gqlgen/graphql"
	"github.com/skinnykaen/rpa_clone/internal/consts"
	"github.com/skinnykaen/rpa_clone/internal/models"
	"github.com/skinnykaen/rpa_clone/pkg/utils"
	"github.com/vektah/gqlparser/v2/gqlerror"
	"log"
	"net/http"
)

func HasRole(errLogger *log.Logger) func(ctx context.Context, obj interface{}, next graphql.Resolver, roles []*models.Role) (interface{}, error) {
	return func(ctx context.Context, obj interface{}, next graphql.Resolver, roles []*models.Role) (interface{}, error) {
		clientRole := ctx.Value(consts.KeyRole)
		if !utils.DoesHaveRole(clientRole.(models.Role), roles) {
			errLogger.Printf("%s", consts.ErrAccessDenied)
			return nil, &gqlerror.Error{
				Extensions: map[string]interface{}{
					"err": utils.ResponseError{
						Code:    http.StatusForbidden,
						Message: consts.ErrAccessDenied,
					},
				},
			}
		}
		return next(ctx)
	}
}
