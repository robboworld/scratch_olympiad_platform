package resolvers

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.
// Code generated by github.com/99designs/gqlgen version v0.17.33

import (
	"context"

	"github.com/robboworld/scratch_olympiad_platform/internal/models"
	"github.com/vektah/gqlparser/v2/gqlerror"
)

// GetAllCountries is the resolver for the GetAllCountries field.
func (r *queryResolver) GetAllCountries(ctx context.Context, page *int, pageSize *int) (*models.CountryHTTPList, error) {
	countries, countRows, err := r.countryService.GetAllCountries(page, pageSize)
	if err != nil {
		r.loggers.Err.Printf("%s", err.Error())
		return nil, &gqlerror.Error{
			Extensions: map[string]interface{}{
				"err": err,
			},
		}
	}
	return &models.CountryHTTPList{
		Countries: models.FromCountriesCore(countries),
		CountRows: int(countRows),
	}, nil
}