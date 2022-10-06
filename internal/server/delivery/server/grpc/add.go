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

func (l LocationService) AddLocation(ctx context.Context, in *generated.AddLocationRequest) (*empty.Empty, error) {
	var description *string
	if in.Location.Description != nil {
		str := in.Location.Description.String()
		description = &str
	}

	loc := &locDomain.Location{
		Caption:     in.Location.Caption,
		Description: description,
	}

	err := l.uc.Add(ctx, loc)
	if err != nil {
		if errors.Is(err, locErrors.ErrLocationInvalid) {
			return &empty.Empty{}, ErrLocationInvalid
		}
		return &empty.Empty{}, fmt.Errorf("%w - %s", ErrLocationUnknownError, err)
	}

	return &empty.Empty{}, nil
}
