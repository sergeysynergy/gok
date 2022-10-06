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

func (l LocationService) DelLocation(ctx context.Context, in *generated.DelLocationRequest) (*empty.Empty, error) {
	err := l.uc.Del(ctx, locDomain.LocationInd(in.Ind))
	if err != nil {
		if errors.Is(err, locErrors.ErrLocationNotFound) {
			return nil, ErrLocationNotFound
		}
		return nil, fmt.Errorf("%w - %s", ErrLocationUnknownError, err)
	}

	return &empty.Empty{}, nil
}
