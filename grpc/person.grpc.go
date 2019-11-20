package grpc

import (
	"context"
	"database/sql"
	"log"

	"github.com/lffranca/trophy-count-model/model"
)

type personGRPC struct {
	model.PersonServer
	DB *sql.DB
}

func (pro *personGRPC) GetByHash(ctx context.Context, req *model.HashRequest) (*model.HashResponse, error) {
	person := struct {
		ID        sql.NullInt64
		Hash      sql.NullString
		FirstName sql.NullString
		LastName  sql.NullString
	}{}

	errPerson := pro.DB.QueryRowContext(
		ctx,
		`
			select
				u.id,
				u.hash,
				u.first_name,
				u.last_name
			from access.user u
			where u.hash = $1
			limit 1;
		`,
		req.GetHash()).Scan(&person.ID, &person.Hash, &person.FirstName, &person.LastName)
	if errPerson != nil {
		log.Println(errPerson)
		return nil, errPerson
	}

	return &model.HashResponse{
		Id:        person.ID.Int64,
		Hash:      person.Hash.String,
		FirstName: person.FirstName.String,
		LastName:  person.LastName.String}, nil
}
