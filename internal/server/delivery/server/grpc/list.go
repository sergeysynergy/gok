package grpc

import (
	generated "back-mfsb/pb/location"
	"context"
	"fmt"
	"github.com/golang/protobuf/ptypes/wrappers"
)

func (l LocationService) ListLocations(ctx context.Context, in *generated.ListLocationsRequest) (*generated.ListLocationsResponse, error) {
	list, err := l.repo.SortByInd(ctx, int(in.Offset), int(in.Limit))
	if err != nil {
		return nil, fmt.Errorf("%w - %s", ErrLocationsList, err)
	}

	listPB := make([]*generated.Location, 0, len(list))
	for _, v := range list {
		description := new(wrappers.StringValue)
		if v.Description != nil {
			description.Value = *v.Description
		}
		listPB = append(listPB, &generated.Location{
			Ind:         int32(v.Ind),
			Caption:     v.Caption,
			Description: description,
		})
	}

	return &generated.ListLocationsResponse{
		Count:     int32(len(list)),
		Locations: listPB,
	}, nil
}
