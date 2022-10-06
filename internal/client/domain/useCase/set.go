package useCase

import (
	"context"
	"fmt"

	locDomain "back-mfsb/service/location/internal/domain/entity"
)

func (l *LocationUseCase) Set(ctx context.Context, loc *locDomain.Location) error {
	errPrefix := "useCase Location.Set"

	err := l.repo.Update(ctx, loc)
	if err != nil {
		return fmt.Errorf("%s - %w", errPrefix, err)
	}

	return nil
}
