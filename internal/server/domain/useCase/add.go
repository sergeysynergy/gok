package useCase

import (
	"context"
	"fmt"

	locDomain "back-mfsb/service/location/internal/domain/entity"
)

func (l *LocationUseCase) Add(ctx context.Context, loc *locDomain.Location) (err error) {
	errPrefix := "useCase Location.Add"

	ind, err := l.repo.Create(ctx, loc)
	if err != nil {
		return fmt.Errorf("%s - %w", errPrefix, err)
	}

	msg := fmt.Sprintf("Location with ID %d successfully created", *ind)
	l.lg.Debug(msg)

	return nil
}
