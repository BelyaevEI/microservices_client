package user

import (
	"context"

	"github.com/BelyaevEI/microservices_auth/internal/client/db"
	sq "github.com/Masterminds/squirrel"
)

func (r *repo) DeleteUserByID(ctx context.Context, id int64) error {
	builderDelete := sq.Delete(tableName).
		PlaceholderFormat(sq.Dollar).
		Where(sq.Eq{idColumn: id})

	query, args, err := builderDelete.ToSql()
	if err != nil {
		return err
	}
	q := db.Query{
		Name:     "user_repository.Delete",
		QueryRaw: query,
	}
	_, err = r.db.DB().ExecContext(ctx, q, args...)
	if err != nil {
		return err
	}
	return nil
}