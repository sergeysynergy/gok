package memory

import (
	locErrors "back-mfsb/service/location/internal/errors"
	"context"

	locDomain "back-mfsb/service/location/internal/domain/entity"
)

func (r *Repo) Update(_ context.Context, loc *locDomain.Location) error {
	r.byIndMu.Lock()
	defer r.byIndMu.Unlock()

	oldLoc, ok := r.byInd[loc.Ind]
	if !ok {
		return locErrors.ErrLocationNotFound
	}

	*oldLoc = *loc

	return nil
}
