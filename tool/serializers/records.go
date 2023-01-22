package serializers

import (
	gokConsts "github.com/sergeysynergy/gok/internal/consts"
	"github.com/sergeysynergy/gok/internal/entity"
	pb "github.com/sergeysynergy/gok/proto"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func RecordPBToEntity(in *pb.Record) *entity.Record {
	rec := &entity.Record{
		ID:          entity.RecordID(in.Id),
		Head:        in.Head,
		BranchID:    entity.BranchID(in.BranchID),
		Description: entity.StringField(in.Description),
		Type:        gokConsts.RecordType(in.Type),
		UpdatedAt:   in.UpdatedAt.AsTime(),
	}
	switch in.Type {
	case string(gokConsts.TEXT):
		rec.Extension = &entity.Text{
			Text: entity.StringField(in.Text.Text),
		}
	case string(gokConsts.PASS):
		rec.Extension = &entity.Pass{
			Login:    entity.StringField(in.Pass.Login),
			Password: entity.StringField(in.Pass.Password),
		}
	case string(gokConsts.CARD):
		rec.Extension = &entity.Card{
			Number:  entity.NumberField(in.Card.Number),
			Code:    entity.NumberField(in.Card.Code),
			Expired: entity.StringField(in.Card.Expired),
			Owner:   entity.StringField(in.Card.Owner),
		}
	case string(gokConsts.FILE):
		rec.Extension = &entity.File{
			File: in.File.File,
		}
	}

	return rec
}

func RecordsPBToEntity(in []*pb.Record) []*entity.Record {
	recs := make([]*entity.Record, 0, len(in))
	for _, v := range in {
		recs = append(recs, RecordPBToEntity(v))
	}

	return recs
}

func RecordEntityToPB(in *entity.Record) *pb.Record {
	recPB := &pb.Record{
		Id:          string(in.ID),
		Head:        in.Head,
		BranchID:    uint64(in.BranchID),
		Description: string(in.Description),
		Type:        string(in.Type),
		UpdatedAt:   timestamppb.New(in.UpdatedAt),
	}
	switch in.Type {
	case gokConsts.TEXT:
		recPB.Text = &pb.Text{
			Text: string(in.Extension.(*entity.Text).Text),
		}
	case gokConsts.PASS:
		recPB.Pass = &pb.Pass{
			Login:    string(in.Extension.(*entity.Pass).Login),
			Password: string(in.Extension.(*entity.Pass).Password),
		}
	case gokConsts.CARD:
		recPB.Card = &pb.Card{
			Number:  uint64(in.Extension.(*entity.Card).Number),
			Code:    uint64(in.Extension.(*entity.Card).Code),
			Expired: string(in.Extension.(*entity.Card).Expired),
			Owner:   string(in.Extension.(*entity.Card).Owner),
		}
	case gokConsts.FILE:
		recPB.File = &pb.File{
			File: in.Extension.(*entity.File).File,
		}
	}

	return recPB
}

func RecordsEntityToPB(recs []*entity.Record) []*pb.Record {
	recsPB := make([]*pb.Record, 0, len(recs))
	for _, v := range recs {
		recsPB = append(recsPB, RecordEntityToPB(v))
	}

	return recsPB
}
