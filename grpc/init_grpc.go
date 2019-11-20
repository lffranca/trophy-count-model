package grpc

import (
	"database/sql"
	"net"

	"github.com/lffranca/trophy-count-model/model"
	"google.golang.org/grpc"
)

const (
	port = ":8000"
)

// InitGRPC InitGRPC
func InitGRPC(db *sql.DB) error {
	lis, err := net.Listen("tcp", port)
	if err != nil {
		return err
	}

	s := grpc.NewServer()

	model.RegisterTransactionServer(s, &transactionGRPC{DB: db})

	if err := s.Serve(lis); err != nil {
		return err
	}

	return nil
}
