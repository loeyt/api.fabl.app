package service

import (
	"context"

	"api.fabl.app/internal/repository"
	"api.fabl.app/internal/session"
	pb "api.fabl.app/v1"
	"github.com/oklog/ulid/v2"
)

type itemServiceServer struct {
	repo repository.ItemRepository

	pb.UnimplementedItemServiceServer
}

func (s *itemServiceServer) Export(ctx context.Context, in *pb.ExportRequest) (*pb.ExportResponse, error) {
	accountID, err := session.Account(ctx)
	if err != nil {
		return nil, err
	}
	id, err := ulid.Parse(in.Id)
	if err != nil {
		return nil, err
	}
	item, err := s.repo.Get(ctx, accountID, id)
	if err != nil {
		return nil, err
	}
	importString, err := item.Export()
	if err != nil {
		return nil, err
	}
	return &pb.ExportResponse{
		ImportString: importString,
	}, nil
}

func (s *itemServiceServer) Get(ctx context.Context, in *pb.GetRequest) (*pb.GetResponse, error) {
	accountID, err := session.Account(ctx)
	if err != nil {
		return nil, err
	}
	id, err := ulid.Parse(in.Id)
	if err != nil {
		return nil, err
	}
	item, err := s.repo.Get(ctx, accountID, id)
	if err != nil {
		return nil, err
	}
	return &pb.GetResponse{
		Data: item.Data,
	}, nil
}

func (s *itemServiceServer) Import(ctx context.Context, in *pb.ImportRequest) (*pb.ImportResponse, error) {
	accountID, err := session.Account(ctx)
	if err != nil {
		return nil, err
	}
	item := &repository.Item{TimeMs: in.TimeMs}
	err = item.Import(in.ImportString)
	if err != nil {
		return nil, err
	}
	// TODO: more validation here, see history for example.
	err = s.repo.Create(ctx, accountID, item)
	if err != nil {
		return nil, err
	}
	return &pb.ImportResponse{
		Id: item.ULID.String(),
	}, nil
}

func (s *itemServiceServer) List(ctx context.Context, in *pb.ListRequest) (*pb.ListResponse, error) {
	accountID, err := session.Account(ctx)
	if err != nil {
		// TODO: when public items exist: change this
		return &pb.ListResponse{
			Items: []*pb.ListResponse_Item{},
		}, nil
	}
	items, err := s.repo.List(ctx, accountID)
	if err != nil {
		return nil, err
	}
	pbItems := make([]*pb.ListResponse_Item, len(items))
	for i, item := range items {
		pbItems[i] = &pb.ListResponse_Item{
			Id:  item.ULID.String(),
			Sum: item.Sum256[:],
		}
	}
	return &pb.ListResponse{
		Items: pbItems,
	}, nil
}

// NewItemServiceServer initializes an ItemServiceServer.
func NewItemServiceServer(repo repository.ItemRepository) pb.ItemServiceServer {
	return &itemServiceServer{
		repo: repo,
	}
}
