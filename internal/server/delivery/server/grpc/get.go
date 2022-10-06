package grpc

import (
	"context"
	"errors"
	"fmt"

	generated "back-mfsb/pb/location"
	locDomain "back-mfsb/service/location/internal/domain/entity"
	locErrors "back-mfsb/service/location/internal/errors"
)

func (l LocationService) GetLocation(ctx context.Context, in *generated.GetLocationRequest) (*generated.GetLocationResponse, error) {
	loc, err := l.uc.Get(ctx, locDomain.LocationInd(in.Ind))
	if err != nil {
		if errors.Is(err, locErrors.ErrLocationNotFound) {
			return nil, ErrLocationNotFound
		}
		return nil, fmt.Errorf("%w - %s", ErrLocationUnknownError, err)
	}

	return &generated.GetLocationResponse{
		Location: &generated.Location{
			Ind:         int32(loc.Ind),
			Caption:     loc.Caption,
			Description: l.StringValue(loc.Description),
		},
	}, nil
}
