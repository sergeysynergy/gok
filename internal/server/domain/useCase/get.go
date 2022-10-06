package useCase

import (
	"context"
	"fmt"

	locDomain "back-mfsb/service/location/internal/domain/entity"
)

func (l *LocationUseCase) Get(ctx context.Context, locInd locDomain.LocationInd) (*locDomain.Location, error) {
	errPrefix := "useCase Location.Get"

	loc, err := l.repo.Read(ctx, locInd)
	if err != nil {
		return nil, fmt.Errorf("%s - %w", errPrefix, err)
	}

	return loc, nil
}
