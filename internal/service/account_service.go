package service

import (
	"context"

	"api.fabl.app/internal/repository"
	"api.fabl.app/internal/session"
	pb "api.fabl.app/v1"
	"github.com/google/uuid"
	grpc_auth "github.com/grpc-ecosystem/go-grpc-middleware/auth"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type accountServiceServer struct {
	repo repository.AccountRepository

	pb.UnimplementedAccountServiceServer
}

func (s *accountServiceServer) CurrentAccount(ctx context.Context, in *pb.CurrentAccountRequest) (*pb.CurrentAccountResponse, error) {
	var (
		account *repository.Account
		err     error
	)

	id, err := session.Account(ctx)

	if err != nil {
		token, err := grpc_auth.AuthFromMD(ctx, "bearer")
		if err != nil {
			// TODO: maybe better error here?
			return nil, err
		}
		account, err = s.repo.FromToken(ctx, token)
		if err != nil {
			// TODO: proper unauthenticated?
			return nil, err
		}
	} else {
		account, err = s.repo.Get(ctx, id)
		if err != nil {
			// TODO: clear session? better error?
			return nil, err
		}
	}
	return &pb.CurrentAccountResponse{
		Account: &pb.Account{
			Id:       account.ID.String(),
			Nickname: account.Nickname,
		},
	}, nil
}

func (s *accountServiceServer) Login(ctx context.Context, in *pb.LoginRequest) (*pb.LoginResponse, error) {
	id, err := uuid.Parse(in.Id)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "bad credentials")
	}
	acc, err := s.repo.Get(ctx, id)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "bad credentials")
	}
	if acc.CheckPassword([]byte(in.Password)) != nil {
		return nil, status.Error(codes.InvalidArgument, "bad credentials")
	}
	err = session.Login(ctx, acc.ID)
	if err != nil {
		return nil, err
	}
	return &pb.LoginResponse{}, nil
}

func (s *accountServiceServer) Logout(ctx context.Context, in *pb.LogoutRequest) (*pb.LogoutResponse, error) {
	err := session.Logout(ctx)
	if err != nil {
		return nil, err
	}
	return &pb.LogoutResponse{}, nil
}

// NewAccountServiceServer initializes an AccountServiceServer.
func NewAccountServiceServer(store repository.AccountRepository) pb.AccountServiceServer {
	return &accountServiceServer{
		repo: store,
	}
}
