package memory

import (
	locDomain "back-mfsb/service/location/internal/domain/entity"
	"context"
)

// SortByInd Возвращает копию слайса Местоположений потокозащищённым методом:
// offset - смещение выборки;
// limit - кол-во элементов.
func (r *Repo) SortByInd(_ context.Context, offset, limit int) (locDomain.SortByInd, error) {
	r.byIndMu.RLock()
	defer r.byIndMu.RUnlock()

	end := offset + limit
	if end > len(r.sortByInd) {
		end = len(r.sortByInd)
	}
	if offset >= end {
		offset = 0
		end = 0
	}

	list := make(locDomain.SortByInd, 0, end-offset)
	for _, v := range r.sortByInd[offset:end] {
		list = append(list, v)
	}

	return list, nil
}
