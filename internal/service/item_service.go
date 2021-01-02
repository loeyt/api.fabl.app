package service

import (
	"context"

	"github.com/oklog/ulid/v2"
	pb "loe.yt/factorio-blueprints/internal/pb/factorio_blueprints/v1"
)

type itemServiceServer struct {
	store ItemStore

	pb.UnimplementedItemServiceServer
}

func (s *itemServiceServer) Export(ctx context.Context, in *pb.ExportRequest) (*pb.ExportResponse, error) {
	id, err := ulid.Parse(in.Id)
	if err != nil {
		return nil, err
	}
	data, err := s.store.GetData(ctx, id)
	if err != nil {
		return nil, err
	}
	return &pb.ExportResponse{
		ImportString: string(data),
	}, nil
}

func (s *itemServiceServer) Import(ctx context.Context, in *pb.ImportRequest) (*pb.ImportResponse, error) {
	data, err := extractItemData(in.ImportString)
	if err != nil {
		return nil, err
	}
	// TODO: more validation here, see history for example.
	item, err := s.store.Create(ctx, data, in.TimeMs)
	if err != nil {
		return nil, err
	}
	return &pb.ImportResponse{
		Id: item.ULID.String(),
	}, nil
}

func (s *itemServiceServer) List(ctx context.Context, in *pb.ListRequest) (*pb.ListResponse, error) {
	items, err := s.store.List(ctx)
	if err != nil {
		return nil, err
	}
	pbItems := make([]*pb.ListResponse_Item, len(items))
	for i, item := range items {
		pbItems[i] = &pb.ListResponse_Item{
			Id:  item.ULID.String(),
			Sum: item.Sum[:],
		}
	}
	return &pb.ListResponse{
		Items: pbItems,
	}, nil
}

// NewItemServiceServer initializes an ItemServiceServer.
func NewItemServiceServer(store ItemStore) pb.ItemServiceServer {
	return &itemServiceServer{
		store: store,
	}
}
