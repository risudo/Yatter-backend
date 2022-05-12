package dao

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"yatter-backend-go/app/domain/object"
	"yatter-backend-go/app/domain/repository"

	"github.com/jmoiron/sqlx"
)

type (
	// Implementation for repository.Account
	relation struct {
		db *sqlx.DB
	}
)

func NewRelation(db *sqlx.DB) repository.Relation {
	return &relation{db: db}
}

func (r *relation) Follow(ctx context.Context, loginID object.AccountID, targetID object.AccountID) error {
	const query = "INSERT INTO relation (following_id, follower_id) VALUES(?, ?)"

	_, err := r.db.ExecContext(ctx, query, loginID, targetID)
	if err != nil {
		return fmt.Errorf("%w", err)
	}
	return nil
}

func (r *relation) IsFollowing(ctx context.Context, accountID object.AccountID, targetID object.AccountID) (bool, error) {
	const query = "SELECT EXISTS(SELECT * FROM relation WHERE following_id = ? AND follower_id = ?) AS existing"

	ex := struct {
		Exist bool `db:"existing"`
	}{}
	err := r.db.QueryRowxContext(ctx, query, accountID, targetID).StructScan(&ex)
	if err != nil {
		return false, fmt.Errorf("%w", err)
	}
	return ex.Exist, nil
}

func (r *relation) Following(ctx context.Context, id object.AccountID, p object.Parameters) ([]object.Account, error) {
	var entity []object.Account
	const query = `
SELECT
	account.id,
	account.username,
	account.create_at,
    CASE
		WHEN NOT EXISTS (
			SELECT
				*
			FROM
				relation AS r
				INNER JOIN account AS a ON r.following_id = a.id
			WHERE
				a.id = account.id
		) THEN 0
		ELSE (
			SELECT
				COUNT(*)
			FROM
				relation
			WHERE
				following_id = (
					SELECT
						a.id
					from
						account a
					WHERE
						a.id = account.id
				)
			GROUP BY
				following_id
		)
	END AS "followingcount",
	CASE
		WHEN NOT EXISTS (
			SELECT
				*
			FROM
				relation AS r
				INNER JOIN account AS a ON r.follower_id = a.id
			WHERE
				account.id = a.id
		) THEN 0
		ELSE (
			SELECT
				COUNT(*)
			FROM
				relation
			WHERE
				follower_id = (
					SELECT
						a.id
					FROM
						account a
					WHERE
						account.id = a.id
				)
			GROUP BY
				follower_id
		)
	END AS "followerscount"
FROM
	account
	JOIN relation ON account.id = relation.follower_id
WHERE
	relation.following_id = ?
ORDER BY
	account.id
LIMIT
	?
	`

	err := r.db.SelectContext(ctx, &entity, query, id, p.Limit)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("%w", err)
	}
	return entity, nil
}

func (r *relation) Followers(ctx context.Context, id object.AccountID, p object.Parameters) ([]object.Account, error) {
	var entity []object.Account
	const query = `
SELECT
	account.id,
	account.username,
	account.create_at,
    CASE
		WHEN NOT EXISTS (
			SELECT
				*
			FROM
				relation AS r
				INNER JOIN account AS a ON r.following_id = a.id
			WHERE
				a.id = account.id
		) THEN 0
		ELSE (
			SELECT
				COUNT(*)
			FROM
				relation
			WHERE
				following_id = (
					SELECT
						a.id
					from
						account a
					WHERE
						a.id = account.id
				)
			GROUP BY
				following_id
		)
	END AS "followingcount",
	CASE
		WHEN NOT EXISTS (
			SELECT
				*
			FROM
				relation AS r
				INNER JOIN account AS a ON r.follower_id = a.id
			WHERE
				account.id = a.id
		) THEN 0
		ELSE (
			SELECT
				COUNT(*)
			FROM
				relation
			WHERE
				follower_id = (
					SELECT
						a.id
					FROM
						account a
					WHERE
						account.id = a.id
				)
			GROUP BY
				follower_id
		)
	END AS "followerscount"
FROM
	account
	JOIN relation ON account.id = relation.following_id
WHERE
	relation.follower_id = ?
	AND account.id < ?
	AND account.id > ?
ORDER BY
	account.id
LIMIT
	?
`

	err := r.db.SelectContext(ctx, &entity, query, id, p.MaxID, p.SinceID, p.Limit)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("%w", err)
	}
	return entity, nil
}

func (r *relation) Unfollow(ctx context.Context, loginID object.AccountID, targetID object.AccountID) error {
	const query = "DELETE FROM relation WHERE following_id = ? AND follower_id = ?"

	_, err := r.db.ExecContext(ctx, query, loginID, targetID)
	if err != nil {
		return fmt.Errorf("%w", err)
	}
	return nil
}
