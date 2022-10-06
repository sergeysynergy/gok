package useCase

import (
	"context"
	"fmt"

	locDomain "back-mfsb/service/location/internal/domain/entity"
)

func (l *LocationUseCase) Del(ctx context.Context, locInd locDomain.LocationInd) error {
	errPrefix := "useCase Location.Del"

	err := l.repo.Delete(ctx, locInd)
	if err != nil {
		return fmt.Errorf("%s - %w", errPrefix, err)
	}

	return nil
}
