package useCase

import (
	"context"
	"fmt"

	locDomain "back-mfsb/service/location/internal/domain/entity"
)

func (l *LocationUseCase) ByInd(ctx context.Context) (locDomain.ByInd, error) {
	errPrefix := "useCase Location.ByID"

	byInd, err := l.repo.ByInd(ctx)
	if err != nil {
		return nil, fmt.Errorf("%s - %w", errPrefix, err)
	}

	return byInd, nil
}
