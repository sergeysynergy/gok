// Package memory Пакет реализует хранилище в памяти для структуры: дерево объектов контроля.
package memory

import (
	"context"
	"fmt"
	"sync"

	locDomain "back-mfsb/service/location/internal/domain/entity"
	"back-mfsb/service/location/internal/domain/useCase"
)

// Repo Реализация операций CRUD (create, read, update, delete), а так же других методов для работы с репозиторием.
type Repo struct {
	byIndMu   sync.RWMutex
	byInd     locDomain.ByInd
	sortByInd locDomain.SortByInd
}

var _ useCase.Repo = new(Repo)

func New() *Repo {
	r := &Repo{
		byInd:     make(locDomain.ByInd, 0),
		sortByInd: make(locDomain.SortByInd, 0),
	}
	r.DebugInit()

	return r
}

func (r *Repo) DebugInit() {
	fmt.Println("DEBUG MEMORY INIT")
	locs := []locDomain.Location{
		{
			Caption: "One",
		},
		{
			Caption: "Two",
		},
		{
			Caption:     "Three",
			Description: func(s string) *string { return &s }("three for all"),
		},
		{
			Caption: "Four",
		},
	}
	ctx := context.Background()
	for _, v := range locs {
		r.Create(ctx, &v)
	}
}
