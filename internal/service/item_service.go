package service

import (
	"context"

	pb "loe.yt/factorio-blueprints/internal/pb/factorio_blueprints/v1"
)

type itemServiceServer struct {
	pb.UnimplementedItemServiceServer
}

func (s *itemServiceServer) Import(ctx context.Context, in *pb.ImportRequest) (*pb.ImportResponse, error) {
	return &pb.ImportResponse{
		Item: &pb.Item{
			ImportString: in.ImportString,
			Item: &pb.Item_BlueprintBook{
				BlueprintBook: &pb.BlueprintBook{Id: "fake"},
			},
		},
	}, nil
}

// NewItemServiceServer initializes an ItemServiceServer.
func NewItemServiceServer() pb.ItemServiceServer {
	return &itemServiceServer{}
}
