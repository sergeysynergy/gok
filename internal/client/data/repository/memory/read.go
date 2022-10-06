package memory

import (
	locErrors "back-mfsb/service/location/internal/errors"
	"context"

	locDomain "back-mfsb/service/location/internal/domain/entity"
)

func (r *Repo) Read(_ context.Context, locInd locDomain.LocationInd) (*locDomain.Location, error) {
	r.byIndMu.RLock()
	defer r.byIndMu.RUnlock()

	loc, ok := r.byInd[locInd]
	if !ok {
		return nil, locErrors.ErrLocationNotFound
	}

	return loc, nil
}
