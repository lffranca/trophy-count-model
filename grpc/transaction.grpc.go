package grpc

import (
	"context"
	"database/sql"
	"log"

	"github.com/lffranca/trophy-count-model/model"
)

type transactionGRPC struct {
	model.TransactionServer
	DB *sql.DB
}

func (pro *transactionGRPC) SetCollectedCoin(ctx context.Context, req *model.CollectedCoinRequest) (*model.CollectedCoinResponse, error) {
	rowsInsert, errInsert := pro.DB.QueryContext(
		ctx,
		"insert into transactions.collected_coin (user_id, value) values ($1, $2);",
		req.GetUserId(),
		req.GetValue())
	if errInsert != nil {
		log.Println(errInsert)
		return nil, errInsert
	}
	rowsInsert.Close()

	valueAll := sql.NullInt64{}

	errAll := pro.DB.QueryRowContext(
		ctx,
		"select sum(cc.value) from transactions.collected_coin cc where cc.user_id = $1;",
		req.GetUserId()).Scan(&valueAll)
	if errAll != nil {
		log.Println(errAll)
		return nil, errAll
	}

	trophy := struct {
		ID       sql.NullInt64
		Starting sql.NullInt64
		Name     sql.NullString
	}{}

	errTrophy := pro.DB.QueryRowContext(
		ctx,
		`
			select
				tc.id,
				tc.starting,
				tc.name
			from dimension.trophy_coin tc
			where tc.starting <= $1
			order by tc.starting desc
			limit 1;
		`,
		valueAll.Int64).Scan(&trophy.ID, &trophy.Starting, &trophy.Name)
	if errTrophy != nil {
		log.Println(errTrophy)
		return nil, errTrophy
	}

	return &model.CollectedCoinResponse{
		ValueAll:           valueAll.Int64,
		TrophyCoinId:       trophy.ID.Int64,
		TrophyCoinStarting: trophy.Starting.Int64,
		TrophyCoinName:     trophy.Name.String}, nil
}

func (pro *transactionGRPC) SetKilledMonster(ctx context.Context, req *model.KilledMonsterRequest) (*model.KilledMonsterResponse, error) {
	rowsInsert, errInsert := pro.DB.QueryContext(
		ctx,
		"insert into transactions.killed_monster (user_id, monster_id) values ($1, $2);",
		req.GetUserId(),
		req.GetMonsterId())
	if errInsert != nil {
		log.Println(errInsert)
		return nil, errInsert
	}
	rowsInsert.Close()

	valueAll := sql.NullInt64{}

	errAll := pro.DB.QueryRowContext(
		ctx,
		"select count(*) from transactions.killed_monster k where k.user_id = $1 and k.monster_id = $2;",
		req.GetUserId(), req.GetMonsterId()).Scan(&valueAll)
	if errAll != nil {
		log.Println(errAll)
		return nil, errAll
	}

	trophy := struct {
		ID        sql.NullInt64
		Starting  sql.NullInt64
		Name      sql.NullString
		MonsterID sql.NullInt64
	}{}

	errTrophy := pro.DB.QueryRowContext(
		ctx,
		`
			select
				tm.id,
				tm.starting,
				tm.name,
				tm.monster_id
			from dimension.trophy_monster tm
			where tm.starting <= $1
			and tm.monster_id = $2
			order by tm.starting desc
			limit 1;
		`,
		valueAll.Int64, req.GetMonsterId()).Scan(&trophy.ID, &trophy.Starting, &trophy.Name, &trophy.MonsterID)
	if errTrophy != nil {
		log.Println(errTrophy)
		return nil, errTrophy
	}

	return &model.KilledMonsterResponse{
		ValueAll:              valueAll.Int64,
		TrophyId:              trophy.ID.Int64,
		TrophyMonsterId:       trophy.MonsterID.Int64,
		TrophyMonsterStarting: trophy.Starting.Int64,
		TrophyMonsterName:     trophy.Name.String}, nil
}

func (pro *transactionGRPC) SetDeath(ctx context.Context, req *model.DeathRequest) (*model.DeathResponse, error) {
	rowsInsert, errInsert := pro.DB.QueryContext(
		ctx,
		"insert into transactions.deaths (user_id, timestamp) values ($1, current_timestamp);",
		req.GetUserId())
	if errInsert != nil {
		log.Println(errInsert)
		return nil, errInsert
	}
	rowsInsert.Close()

	valueAll := sql.NullInt64{}

	errAll := pro.DB.QueryRowContext(
		ctx,
		"select count(*) from transactions.deaths d where d.user_id = $1;",
		req.GetUserId()).Scan(&valueAll)
	if errAll != nil {
		log.Println(errAll)
		return nil, errAll
	}

	trophy := struct {
		ID       sql.NullInt64
		Starting sql.NullInt64
		Name     sql.NullString
	}{}

	errTrophy := pro.DB.QueryRowContext(
		ctx,
		`
			select
				d.id,
				d.starting,
				d.name
			from dimension.trophy_death d
			where d.starting <= $1
			order by d.starting desc
			limit 1;
		`,
		valueAll.Int64).Scan(&trophy.ID, &trophy.Starting, &trophy.Name)
	if errTrophy != nil {
		log.Println(errTrophy)
		return nil, errTrophy
	}

	return &model.DeathResponse{
		ValueAll:            valueAll.Int64,
		TrophyDeathId:       trophy.ID.Int64,
		TrophyDeathStarting: trophy.Starting.Int64,
		TrophyDeathName:     trophy.Name.String}, nil
}
