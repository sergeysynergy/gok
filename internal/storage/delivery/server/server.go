// Package server contains gRPC API endpoint to work with GoK server side service.
package server

import (
	"context"
	"github.com/sergeysynergy/gok/internal/auth/useCase"

	"github.com/golang/protobuf/ptypes/empty"
	"go.uber.org/zap"

	pb "github.com/sergeysynergy/gok/proto"
)

type StorageServer struct {
	pb.UnimplementedUsersServer
	lg   *zap.Logger
	user useCase.UseCase
}

//var _ generated.PersonsServer = new(PersonsService)

func New(
	logger *zap.Logger,
	uc useCase.UseCase,
) *StorageServer {
	return &StorageServer{
		lg:   logger,
		user: uc,
	}
}

func (s StorageServer) AddUser(ctx context.Context, in *pb.AddUserRequest) (*empty.Empty, error) {
	//usr := &entity.User{
	//	Login: in.Login,
	//}

	//err := p.uc.Add(ctx, psn)
	//if err != nil {
	//	if errors.Is(err, psnErrors.ErrPersonInvalidArgument) {
	//		return &empty.Empty{}, ErrPersonInvalidArgument
	//	}
	//	if errors.Is(err, psnErrors.ErrPersonAlreadyExists) {
	//		return &empty.Empty{}, ErrPersonAlreadyExists
	//	}
	//	return &empty.Empty{}, fmt.Errorf("%w - %s", ErrPersonUnknownError, err)
	//}

	return &empty.Empty{}, nil
}

func (s StorageServer) ListUsers(ctx context.Context, in *pb.ListUsersRequest) (*pb.ListUsersResponse, error) {
	//list, err := p.repo.List(ctx, int(in.Offset), int(in.Limit))
	//if err != nil {
	//	return nil, fmt.Errorf("%w - %s", ErrPersonsList, err)
	//}
	//
	//personsPB := make([]*generated.Person, 0, len(list.Persons))
	//for _, v := range list.Persons {
	//	patronymic := new(wrappers.StringValue)
	//	if v.Patronymic.Valid {
	//		patronymic.Value = v.Patronymic.String
	//	}
	//	personsPB = append(personsPB, &generated.Person{
	//		Ind:        int32(v.Ind),
	//		Name:       v.Name,
	//		Surname:    v.Surname,
	//		Patronymic: patronymic,
	//		Email:      v.Email,
	//		Phone:      v.Phone,
	//		Active:     v.Active,
	//	})
	//}
	//
	//roleTypesPB := make([]*generated.RoleType, 0, len(list.RoleTypes))
	//for _, v := range list.RoleTypes {
	//	description := new(wrappers.StringValue)
	//	if v.Description.Valid {
	//		description.Value = v.Description.String
	//	}
	//	roleTypesPB = append(roleTypesPB, &generated.RoleType{
	//		Ind:         int32(v.Ind),
	//		Caption:     v.Caption,
	//		Description: description,
	//	})
	//}

	return &pb.ListUsersResponse{
		//	Persons:   personsPB,
		//	RoleTypes: roleTypesPB,
	}, nil
}

//func (p PersonsService) DelPerson(ctx context.Context, in *generated.DelLocationRequest) (*empty.Empty, error) {
//	err := l.uc.Del(ctx, locDomain.LocationInd(in.Ind))
//	if err != nil {
//		if errors.Is(err, locErrors.ErrLocationNotFound) {
//			return nil, ErrLocationNotFound
//		}
//		return nil, fmt.Errorf("%w - %s", ErrLocationUnknownError, err)
//	}
//
//	return &empty.Empty{}, nil
//}

//func (p PersonsService) GetPerson(ctx context.Context, in *generated.GetLocationRequest) (*generated.GetLocationResponse, error) {
//	loc, err := l.uc.Get(ctx, locDomain.LocationInd(in.Ind))
//	if err != nil {
//		if errors.Is(err, locErrors.ErrLocationNotFound) {
//			return nil, ErrLocationNotFound
//		}
//		return nil, fmt.Errorf("%w - %s", ErrLocationUnknownError, err)
//	}
//
//	return &generated.GetLocationResponse{
//		Location: &generated.Location{
//			Ind:         int32(loc.Ind),
//			Caption:     loc.Caption,
//			Description: l.StringValue(loc.Description),
//		},
//	}, nil
//}

//func (p PersonsService) SetPerson(ctx context.Context, in *generated.SetLocationRequest) (*empty.Empty, error) {
//	var description *string
//	if in.Location.Description != nil {
//		description = func(s string) *string { return &s }(in.Location.Description.String())
//	}
//
//	loc := &locDomain.Location{
//		Ind:         locDomain.LocationInd(in.Location.Ind),
//		Caption:     in.Location.Caption,
//		Description: description,
//	}
//	err := l.uc.Set(ctx, loc)
//	if err != nil {
//		if errors.Is(err, locErrors.ErrLocationNotFound) {
//			return nil, ErrLocationNotFound
//		}
//		return nil, fmt.Errorf("%w - %s", ErrLocationUnknownError, err)
//	}
//
//	return &empty.Empty{}, nil
//}
