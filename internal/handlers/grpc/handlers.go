package grpc

import (
	"context"

	pb "github.com/AbramovArseniy/Companies/internal/handlers/grpc/proto"
	db "github.com/AbramovArseniy/Companies/internal/storage/postgres/db"
)

type CompaniesServer struct {
	pb.UnimplementedCompaniesServiceServer

	Storage db.Querier
}

func (s *CompaniesServer) GetTree(context.Context, *pb.GetTreeRequest) (*pb.GetTreeResponse, error) {
	return nil, nil
}

func (s *CompaniesServer) GetHierarchy(_ context.Context, req *pb.GetHierarchyRequest) (*pb.GetHierarchyResponse, error) {
	return nil, nil
}

func (s *CompaniesServer) GetNode(_ context.Context, req *pb.GetNodeRequest) (*pb.GetNodeResponse, error) {
	return nil, nil
}
