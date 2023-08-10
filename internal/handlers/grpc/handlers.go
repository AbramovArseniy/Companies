package grpc

import (
	"context"
	"database/sql"
	"log"

	"github.com/AbramovArseniy/Companies/internal/cfg"
	pb "github.com/AbramovArseniy/Companies/internal/handlers/grpc/proto"
	db "github.com/AbramovArseniy/Companies/internal/storage/postgres/db"
)

type CompaniesServer struct {
	pb.UnimplementedCompaniesServiceServer

	Storage db.Querier
}

func New(cfg *cfg.Config) *CompaniesServer {
	database, err := sql.Open("pgx", cfg.DBAddress)
	if err != nil {
		log.Println("error while opening database:", err)
		return nil
	}
	querier := db.New(database)
	return &CompaniesServer{
		Storage: querier,
	}
}

func (s *CompaniesServer) GetTree(context.Context, *pb.GetTreeRequest) (*pb.GetTreeResponse, error) {
	tree, err := s.Storage.GetAllTree(context.Background())
	if err != nil {
		log.Println("error while getting tree from database:", err)
		return nil, err
	}
	resp := &pb.GetTreeResponse{
		Info: make([]*pb.NodeInfo, len(tree)),
	}
	for i, val := range tree {
		resp.Info[i].Address = val.Address.String
		resp.Info[i].ContactPerson = val.ContactPerson.String
		resp.Info[i].Id = val.ID
		resp.Info[i].Name = val.Name
		resp.Info[i].ParentId = val.ParentID.Int32
		resp.Info[i].PhoneNumber = val.PhoneNumber.String
	}
	return resp, nil
}

func (s *CompaniesServer) GetHierarchy(_ context.Context, req *pb.GetHierarchyRequest) (*pb.GetHierarchyResponse, error) {
	tree, err := s.Storage.GetHierarchy(context.Background(), req.NodeId)
	if err != nil {
		log.Println("error while getting tree from database:", err)
		return nil, err
	}
	resp := &pb.GetHierarchyResponse{
		Info: make([]*pb.NodeInfo, len(tree)),
	}
	for i, val := range tree {
		resp.Info[i].Address = val.Address.String
		resp.Info[i].ContactPerson = val.ContactPerson.String
		resp.Info[i].Id = val.ID.Int32
		resp.Info[i].Name = val.Name.String
		resp.Info[i].ParentId = val.ParentID.Int32
		resp.Info[i].PhoneNumber = val.PhoneNumber.String
	}
	return resp, nil
}

func (s *CompaniesServer) GetNode(_ context.Context, req *pb.GetNodeRequest) (*pb.GetNodeResponse, error) {
	node, err := s.Storage.GetOneNode(context.Background(), req.NodeId)
	if err != nil {
		log.Println("error while getting tree from database:", err)
		return nil, err
	}
	resp := &pb.GetNodeResponse{
		Info: &pb.NodeInfo{
			Address:       node.Address.String,
			ContactPerson: node.ContactPerson.String,
			Id:            node.ID,
			Name:          node.Name,
			ParentId:      node.ParentID.Int32,
		},
	}
	return resp, nil
}
