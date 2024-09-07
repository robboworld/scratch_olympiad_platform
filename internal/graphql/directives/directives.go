package directives

import (
	"context"
	"github.com/99designs/gqlgen/graphql"
	"github.com/robboworld/scratch_olympiad_platform/internal/consts"
	"github.com/robboworld/scratch_olympiad_platform/internal/models"
	"github.com/robboworld/scratch_olympiad_platform/pkg/utils"

	"github.com/vektah/gqlparser/v2/gqlerror"
	"log"
	"net/http"
)

func HasRole(errLogger *log.Logger) func(ctx context.Context, obj interface{}, next graphql.Resolver, roles []*models.Role) (interface{}, error) {
	return func(ctx context.Context, obj interface{}, next graphql.Resolver, roles []*models.Role) (interface{}, error) {
		ginContext, err := utils.GinContextFromContext(ctx)
		if err != nil {
			errLogger.Printf("%s", err.Error())
			return nil, &gqlerror.Error{
				Extensions: map[string]interface{}{
					"err": err,
				},
			}
		}
		clientRole := ginContext.Value(consts.KeyRole)
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
