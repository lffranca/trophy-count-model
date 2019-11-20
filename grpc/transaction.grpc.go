package grpc

import (
	"context"
	"database/sql"

	"github.com/lffranca/trophy-count-model/model"
)

type transactionGRPC struct {
	model.TransactionServer
	DB *sql.DB
}

func (pro *transactionGRPC) SetCollectedCoin(context.Context, *model.CollectedCoinRequest) (*model.CollectedCoinResponse, error) {
	return &model.CollectedCoinResponse{}, nil
}

func (pro *transactionGRPC) SetKilledMonster(context.Context, *model.KilledMonsterRequest) (*model.KilledMonsterResponse, error) {
	return &model.KilledMonsterResponse{}, nil
}

func (pro *transactionGRPC) SetDeath(context.Context, *model.DeathRequest) (*model.DeathResponse, error) {
	return &model.DeathResponse{}, nil
}
