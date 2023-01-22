package useCase

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	"go.uber.org/zap"
	"time"

	gokClient "github.com/sergeysynergy/gok/internal/cli/delivery/client"
	recRepo "github.com/sergeysynergy/gok/internal/data/repository/sql/record"
	"github.com/sergeysynergy/gok/internal/entity"
	gokErrors "github.com/sergeysynergy/gok/internal/errors"
)

type GokUseCase struct {
	lg     *zap.Logger
	ctx    context.Context
	repo   Repo
	client Client
}

func New(logger *zap.Logger, repo *recRepo.Repo, client *gokClient.GokClient) *GokUseCase {
	const (
		defaultTimeout = 120 * time.Second
	)
	ctx, _ := context.WithTimeout(context.Background(), defaultTimeout)

	uc := &GokUseCase{
		ctx:    ctx,
		lg:     logger,
		repo:   repo,
		client: client,
	}

	return uc
}

func (u *GokUseCase) SignIn(usrCLI *entity.CLIUser) (*entity.SignedUser, error) {
	var err error
	defer func() {
		prefix := "GokUseCase.SignIn"
		if err != nil {
			err = fmt.Errorf("%s - %w", prefix, err)
			u.lg.Error(err.Error())
		}
	}()

	usr := &entity.User{
		Login: usrCLI.Login,
	}

	signedUsr, err := u.client.SignIn(u.ctx, usr)
	if err != nil {
		return nil, err
	}

	return signedUsr, nil
}

func (u *GokUseCase) Login(usrCLI *entity.CLIUser) (*entity.SignedUser, error) {
	var err error
	defer func() {
		prefix := "GokUseCase.Login"
		if err != nil {
			err = fmt.Errorf("%s - %w", prefix, err)
			u.lg.Error(err.Error())
		}
	}()

	usr := &entity.User{
		Login: usrCLI.Login,
	}

	signedUsr, err := u.client.Login(u.ctx, usr)
	if err != nil {
		return nil, fmt.Errorf("%w -%s", gokErrors.ErrLoginFailed, err)
	}

	return signedUsr, nil
}

func (u *GokUseCase) Init(token string, localHead uint64) (*entity.Branch, error) {
	var err error
	defer func() {
		prefix := "GokUseCase.Init"
		if err != nil {
			err = fmt.Errorf("%s - %w", prefix, err)
			u.lg.Error(err.Error())
		} else {
			u.lg.Debug(fmt.Sprintf("%s done successfully", prefix))
		}
	}()

	if localHead > 0 {
		return nil, fmt.Errorf("branch already has been initiated - try pull")
	}

	brn, err := u.client.Init(u.ctx, token)
	if err != nil {
		return nil, err
	}
	u.lg.Debug(fmt.Sprintf("GokUseCase.Init - got remote branch ID %d, name %s, head %d", brn.ID, brn.Name, brn.Head))

	locBrn := &entity.Branch{
		ID:   brn.ID,
		Head: localHead,
	}

	if brn.Head > localHead {
		u.lg.Debug("GokUseCase.Init - branch already exists on server, doing force pull to init new local repository")

		freshBrn, freshRecs, errPull := u.client.Pull(u.ctx, token, locBrn)
		if len(freshRecs) == 0 {
			return brn, gokErrors.ErrRecordNotFound
		}
		if errPull != nil {
			return nil, errPull
		}
		if err = u.repo.BulkCreateUpdate(u.ctx, freshRecs); err != nil {
			return nil, fmt.Errorf("%w - %s", gokErrors.ErrPullFailed, err)
		}
		return freshBrn, nil
	}

	return brn, nil
}

func (u *GokUseCase) Push(token string, brn *entity.Branch) (*entity.Branch, error) {
	var err error
	logPrefix := "GokUseCase.Push"
	defer func() {
		if err != nil {
			err = fmt.Errorf("%s - %w", logPrefix, err)
			u.lg.Error(err.Error())
		} else {
			u.lg.Debug(fmt.Sprintf("%s done successfully", logPrefix))
		}
	}()

	recs, err := u.repo.HeadList(u.ctx, brn.ID, brn.Head)
	if err != nil {
		return nil, err
	}

	brn, err = u.client.Push(u.ctx, token, brn, recs)
	if err != nil {
		return nil, err
	}

	u.lg.Debug(fmt.Sprintf("%s got branch: ID %d; name %s; head %d", logPrefix, brn.ID, brn.Name, brn.Head))
	return brn, nil
}

func (u *GokUseCase) Pull(cfg *entity.CLIConf, locBrn *entity.Branch) (*entity.Branch, error) {
	u.lg.Debug("doing GokUseCase.Pull")
	var err error
	defer func() {
		prefix := "GokUseCase.Pull"
		if err != nil {
			err = fmt.Errorf("%s - %w", prefix, err)
			u.lg.Error(err.Error())
		} else {
			u.lg.Debug(fmt.Sprintf("%s done successfully", prefix))
		}
	}()

	freshBrn, freshRecs, err := u.client.Pull(u.ctx, cfg.Token, locBrn)
	if err != nil {
		return nil, err
	}

	// Get records IDs and create map for later merging.
	ids := make([]string, 0, len(freshRecs))
	for _, v := range freshRecs {
		ids = append(ids, string(v.ID))
	}
	u.lg.Debug(fmt.Sprintf("records IDs for merging: %s", ids))

	// Get local records with given ids.
	locRecs, err := u.repo.ByIDsList(u.ctx, ids)
	if err != nil && err != gokErrors.ErrRecordNotFound {
		return nil, fmt.Errorf("%w - %s", gokErrors.ErrPullFailed, err)
	}
	locRecsByID := make(map[entity.RecordID]*entity.Record, len(locRecs))
	for _, v := range locRecs {
		locRecsByID[v.ID] = v
	}

	u.lg.Debug(fmt.Sprintf("local branch header: %d", locBrn.Head))
	u.lg.Debug(fmt.Sprintf("fresh branch header: %d", freshBrn.Head))
	for _, freshRec := range freshRecs {
		// Check if server record already exists locally.
		locRec, ok := locRecsByID[freshRec.ID]
		if ok {
			if locRec.Head > locBrn.Head {
				u.lg.Debug("local record has been changed: have to solve merge conflicts")
				err = u.resolveConflicts(freshBrn, freshRec, locRec)
				if err != nil {
					return nil, err
				}
			}
		}
	}

	// Finally write updated records to repository.
	if err = u.repo.BulkCreateUpdate(u.ctx, freshRecs); err != nil {
		return nil, fmt.Errorf("%w - %s", gokErrors.ErrPullFailed, err)
	}

	u.lg.Debug(fmt.Sprintf("%d records have been merged successfully", len(locRecs)))
	return freshBrn, nil
}

func (u *GokUseCase) resolveConflicts(freshBrn *entity.Branch, freshRec, locRec *entity.Record) error {
	var err error
	defer func() {
		prefix := "GokUseCase.resolveConflicts"
		if err != nil {
			err = fmt.Errorf("%s - %w", prefix, err)
			u.lg.Error(err.Error())
		} else {
			u.lg.Debug(fmt.Sprintf("%s done successfully", prefix))
		}
	}()

	msg := `Oops! We got a merge conflict! Make you chose (1 - is safe unbrained choice):
	1 - Clone updated local record; then replace it with new server version;
	2 - skipp all changes in local record - all changes lost; replace local record with new server version;
	3 - keep local record, ignore server version: set record head to server so record will be included in future push.`
	fmt.Println(msg)

	var choice int8
	_, err = fmt.Scanf("%d", &choice)
	if err != nil {
		return fmt.Errorf("%w - %s", gokErrors.ErrResolveConflict, err)
	}
	switch choice {
	case 1:
		u.lg.Debug("clone updated local record; then replace it with new server version.")
		if err = u.clone(freshBrn, locRec); err != nil {
			return fmt.Errorf("%w - %s", gokErrors.ErrResolveConflict, err)
		}
		*locRec = *freshRec
	case 2:
		u.lg.Debug("skipp all changes in local record - all changes lost; replace local record with new server version.")
	case 3:
		u.lg.Debug("keep local record, ignore server version: set record head to server+1 so record will be included in future push.")
		u.keepLocal(freshBrn, freshRec, locRec)
	default:
		err = fmt.Errorf("invalid choice")
		err = fmt.Errorf("%w - %s", gokErrors.ErrResolveConflict, err)
		return err
	}

	return nil
}

func (u *GokUseCase) clone(freshBrn *entity.Branch, rec *entity.Record) error {
	var err error
	defer func() {
		prefix := "GokUseCase.clone"
		if err != nil {
			err = fmt.Errorf("%s - %w", prefix, err)
			u.lg.Error(err.Error())
		} else {
			u.lg.Debug(fmt.Sprintf("%s done successfully", prefix))
		}
	}()

	clonedRec := &entity.Record{
		ID:          entity.RecordID(uuid.New().String()),
		Head:        freshBrn.Head + 1, // set to server branch head+1, so record will be included in next push
		BranchID:    rec.BranchID,
		Description: rec.Description,
		UpdatedAt:   time.Now(),
		Type:        rec.Type,
		Extension:   rec.Extension,
	}
	err = u.repo.Create(u.ctx, clonedRec)
	if err != nil {
		return fmt.Errorf("%w - %s", gokErrors.ErrCloningRecord, err)
	}

	return nil
}

func (u *GokUseCase) keepLocal(freshBrn *entity.Branch, freshRec, locRec *entity.Record) {
	*freshRec = *locRec
	freshRec.Head = freshBrn.Head + 1 // set to server branch head +1, so record will be included in next push
	freshRec.UpdatedAt = time.Now()
}
