package grpc

import (
	"context"
	"database/sql"
	"log"

	"github.com/lffranca/trophy-count-model/model"
)

type trophyGRPC struct {
	model.TrophyServer
	DB *sql.DB
}

type killedMonsterCount struct {
	MonsterID sql.NullInt64
	Total     sql.NullInt64
}

type trophyMonster struct {
	ID        sql.NullInt64
	Starting  sql.NullInt64
	Name      sql.NullString
	MonsterID sql.NullInt64
}

func (pro *trophyGRPC) GetTrophyCoin(ctx context.Context, req *model.DefaultRequest) (*model.TrophyCoinResponse, error) {
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

	return &model.TrophyCoinResponse{
		ValueAll:           valueAll.Int64,
		TrophyCoinId:       trophy.ID.Int64,
		TrophyCoinStarting: trophy.Starting.Int64,
		TrophyCoinName:     trophy.Name.String}, nil
}

func (pro *trophyGRPC) GetTrophyDeath(ctx context.Context, req *model.DefaultRequest) (*model.TrophyDeathResponse, error) {
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

	return &model.TrophyDeathResponse{
		ValueAll:            valueAll.Int64,
		TrophyDeathId:       trophy.ID.Int64,
		TrophyDeathStarting: trophy.Starting.Int64,
		TrophyDeathName:     trophy.Name.String}, nil
}

func (pro *trophyGRPC) GetTrophyMonster(ctx context.Context, req *model.DefaultRequest) (*model.TrophyMonsterResponse, error) {
	rowsCount, errCount := pro.DB.QueryContext(
		ctx,
		`
			select
				k.monster_id,
				count(*) as total
			from transactions.killed_monster k
			where k.user_id = $1
			group by k.monster_id
			;
		`,
		req.GetUserId())
	if errCount != nil {
		log.Println(errCount)
		return nil, errCount
	}

	count := []killedMonsterCount{}
	countItem := killedMonsterCount{}
	for rowsCount.Next() {
		if err := rowsCount.Scan(
			&countItem.MonsterID,
			&countItem.Total); err != nil {
			log.Println(err)
			continue
		}

		count = append(count, countItem)
	}

	rowsCount.Close()

	killedMonster := []*model.KilledMonster{}

	for _, item := range count {
		itemResult := trophyMonster{}

		errItem := pro.DB.QueryRowContext(
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
			item.Total.Int64, item.MonsterID.Int64).Scan(
			&itemResult.ID,
			&itemResult.Starting,
			&itemResult.Name,
			&itemResult.MonsterID)
		if errItem != nil {
			log.Println(errItem)
			continue
		}

		killedMonster = append(killedMonster, &model.KilledMonster{
			TrophyId:              itemResult.ID.Int64,
			TrophyMonsterId:       itemResult.MonsterID.Int64,
			TrophyMonsterStarting: itemResult.Starting.Int64,
			TrophyMonsterName:     itemResult.Name.String})

	}

	return &model.TrophyMonsterResponse{KilledMonster: killedMonster}, nil
}
