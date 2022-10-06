package memory

import (
	locDomain "back-mfsb/service/location/internal/domain/entity"
	"context"
	"fmt"
)

func (r *Repo) Create(_ context.Context, newLoc *locDomain.Location) (*locDomain.LocationInd, error) {
	loc := *newLoc
	r.byIndMu.Lock()
	defer r.byIndMu.Unlock()

	newInd := locDomain.LocationInd(len(r.sortByInd) + 1)
	loc.Ind = newInd
	fmt.Println(":: ADDING LOC", loc)
	r.sortByInd = append(r.sortByInd, &loc)
	r.byInd[loc.Ind] = &loc

	return &loc.Ind, nil
}
