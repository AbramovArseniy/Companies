// package grpc describes gRPC server's work
package grpc

import (
	"context"
	"fmt"
	"log"

	pb "github.com/AbramovArseniy/Companies/internal/handlers/grpc/proto"
	db "github.com/AbramovArseniy/Companies/internal/storage/postgres/db"
	"github.com/jackc/pgx/v5/pgxpool"
)

// CompaniesServer describes grpc server
type CompaniesServer struct {
	pb.UnimplementedCompaniesServiceServer

	Storage db.Querier
}

// New creates new CompaniesServer from config
func New(dbPool *pgxpool.Pool) (*CompaniesServer, error) {
	dbConn, err := dbPool.Acquire(context.Background())
	if err != nil {
		return nil, fmt.Errorf("error while acquiring database connection: %w", err)
	}
	storage := db.New(dbConn)

	return &CompaniesServer{
		Storage: storage,
	}, nil
}

// GetTree returns information about all the nodes in the tree
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
		resp.Info[i] = &pb.NodeInfo{
			Id:   val.ID,
			Name: val.Name,
		}
		if val.Address.Valid {
			resp.Info[i].Address = val.Address.String
		}
		if val.ContactPerson.Valid {
			resp.Info[i].ContactPerson = val.ContactPerson.String
		}
		resp.Info[i].Id = val.ID
		resp.Info[i].Name = val.Name
		if val.ParentID.Valid {
			resp.Info[i].ParentId = val.ParentID.Int32
		}
		if val.PhoneNumber.Valid {
			resp.Info[i].PhoneNumber = val.PhoneNumber.String
		}
	}
	return resp, nil
}

// GetHierarchy returns information about hierarchy of a node by the node id
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
		resp.Info[i] = &pb.NodeInfo{}
		if val.Address.Valid {
			resp.Info[i].Address = val.Address.String
		}
		if val.ContactPerson.Valid {
			resp.Info[i].ContactPerson = val.ContactPerson.String
		}
		if val.ID.Valid {
			resp.Info[i].Id = val.ID.Int32
		}
		if val.Name.Valid {
			resp.Info[i].Name = val.Name.String
		}
		if val.ParentID.Valid {
			resp.Info[i].ParentId = val.ParentID.Int32
		}
		if val.PhoneNumber.Valid {
			resp.Info[i].PhoneNumber = val.PhoneNumber.String
		}
	}
	return resp, nil
}

// GetNode returns information about one node by the node id
func (s *CompaniesServer) GetNode(_ context.Context, req *pb.GetNodeRequest) (*pb.GetNodeResponse, error) {
	node, err := s.Storage.GetOneNode(context.Background(), req.NodeId)
	if err != nil {
		log.Println("error while getting tree from database:", err)
		return nil, err
	}
	resp := &pb.GetNodeResponse{
		Info: &pb.NodeInfo{
			Id:   node.ID,
			Name: node.Name,
		},
	}

	if node.Address.Valid {
		resp.Info.Address = node.Address.String
	}
	if node.ContactPerson.Valid {
		resp.Info.ContactPerson = node.ContactPerson.String
	}
	if node.ParentID.Valid {
		resp.Info.ParentId = node.ParentID.Int32
	}
	if node.PhoneNumber.Valid {
		resp.Info.PhoneNumber = node.PhoneNumber.String
	}
	return resp, nil
}
