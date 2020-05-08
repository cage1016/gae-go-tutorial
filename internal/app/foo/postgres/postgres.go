package postgres

import (
	"context"

	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
	"github.com/gomurphyx/sqlx"

	"github.com/cage1016/gae-lab-001/internal/app/foo/model"
	"github.com/cage1016/gae-lab-001/internal/pkg/errors"
	"github.com/cage1016/gae-lab-001/internal/pkg/responses"
)

var _ model.FooRepository = (*fooRepository)(nil)

var (
	ErrInsertOrUpdateToFeedbackDB = errors.New("insert or update DB failed")
)

type fooRepository struct {
	db  *sqlx.DB
	log log.Logger
}

// New instantiates a PostgreSQL implementation of givenEmail
// repository.
func New(db *sqlx.DB, log log.Logger) model.FooRepository {
	return &fooRepository{db, log}
}

func (f fooRepository) Insert(ctx context.Context, foo model.Foo) (string, error) {
	q := `INSERT INTO foos (id, value) VALUES (:id, :value);`
	if _, err := f.db.NamedExecContext(ctx, q, foo); err != nil {
		level.Error(f.log).Log("method", "f.db.NamedExecContext", "err", err)
		return "", errors.Wrap(ErrInsertOrUpdateToFeedbackDB, err)
	}
	return foo.ID, nil
}

func (f fooRepository) RetrieveAll(ctx context.Context, offset uint64, limit uint64) (model.FooItemPage, error) {
	items := []model.Foo{}
	q := `SELECT * FROM foos order by created_at desc limit $1 offset $2`
	if err := f.db.SelectContext(ctx, &items, q, limit, offset); err != nil {
		level.Error(f.log).Log("method", "f.db.SelectContext", "sql", q, "limit", limit, "offset", offset)
		return model.FooItemPage{Items: []model.Foo{}}, err
	}

	total, err := total(ctx, f.db, `select count(*) from foos`, map[string]interface{}{})
	if err != nil {
		return model.FooItemPage{Items: []model.Foo{}}, err
	}

	return model.FooItemPage{
		Items: items,
		Paging: responses.Paging{
			Total:  total,
			Limit:  limit,
			Offset: offset,
		},
	}, nil
}

func total(ctx context.Context, db *sqlx.DB, query string, params map[string]interface{}) (uint64, error) {
	rows, err := db.NamedQueryContext(ctx, query, params)
	if err != nil {
		return 0, err
	}

	total := uint64(0)
	if rows.Next() {
		if err := rows.Scan(&total); err != nil {
			return 0, err
		}
	}

	return total, nil
}
