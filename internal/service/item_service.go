package service

import (
	"context"

	"google.golang.org/protobuf/encoding/protojson"
	pb "loe.yt/factorio-blueprints/internal/pb/factorio_blueprints/v1"
)

type itemServiceServer struct {
	pb.UnimplementedItemServiceServer
}

func (s *itemServiceServer) Import(ctx context.Context, in *pb.ImportRequest) (*pb.ImportResponse, error) {
	b, err := extractImportString(in.ImportString)
	if err != nil {
		return nil, err
	}
	m := new(pb.Item)
	err = protojson.UnmarshalOptions{
		DiscardUnknown: true,
	}.Unmarshal(b, m)
	if err != nil {
		return nil, err
	}
	return &pb.ImportResponse{
		Item: &pb.Item{
			ImportString: in.ImportString,
			Item:         m.Item,
		},
	}, nil
}

// NewItemServiceServer initializes an ItemServiceServer.
func NewItemServiceServer() pb.ItemServiceServer {
	return &itemServiceServer{}
}
