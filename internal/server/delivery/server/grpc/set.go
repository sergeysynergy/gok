package grpc

import (
	"context"
	"errors"
	"fmt"
	"github.com/golang/protobuf/ptypes/empty"

	generated "back-mfsb/pb/location"
	locDomain "back-mfsb/service/location/internal/domain/entity"
	locErrors "back-mfsb/service/location/internal/errors"
)

func (l LocationService) SetLocation(ctx context.Context, in *generated.SetLocationRequest) (*empty.Empty, error) {
	var description *string
	if in.Location.Description != nil {
		description = func(s string) *string { return &s }(in.Location.Description.String())
	}

	loc := &locDomain.Location{
		Ind:         locDomain.LocationInd(in.Location.Ind),
		Caption:     in.Location.Caption,
		Description: description,
	}
	err := l.uc.Set(ctx, loc)
	if err != nil {
		if errors.Is(err, locErrors.ErrLocationNotFound) {
			return nil, ErrLocationNotFound
		}
		return nil, fmt.Errorf("%w - %s", ErrLocationUnknownError, err)
	}

	return &empty.Empty{}, nil
}
