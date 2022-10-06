package memory

import (
	locErrors "back-mfsb/service/location/internal/errors"
	"context"

	locDomain "back-mfsb/service/location/internal/domain/entity"
)

func (r *Repo) Delete(_ context.Context, locInd locDomain.LocationInd) error {
	r.byIndMu.Lock()
	defer r.byIndMu.Unlock()

	_, ok := r.byInd[locInd]
	if !ok {
		return locErrors.ErrLocationNotFound
	}

	delete(r.byInd, locInd)

	for k, v := range r.sortByInd {
		if locInd == v.Ind {
			r.sortByInd = append(r.sortByInd[:k], r.sortByInd[k+1:]...)
			break
		}
	}

	return nil
}
