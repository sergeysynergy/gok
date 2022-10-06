package memory

import (
	"context"

	locDomain "back-mfsb/service/location/internal/domain/entity"
)

func (r *Repo) ByInd(_ context.Context) (locDomain.ByInd, error) {
	byInd := make(locDomain.ByInd, len(r.byInd))

	// Выполняем потокозащищённое копирование мапы, иначе будет возвращён указатель на мапу хранилища.
	r.byIndMu.RLock()
	defer r.byIndMu.RUnlock()
	for k, v := range r.byInd {
		byInd[k] = v
	}

	return byInd, nil
}
