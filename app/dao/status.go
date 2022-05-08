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
	// Implementation for repository.Status
	status struct {
		db *sqlx.DB
	}
)

// Create status repository
func NewStatus(db *sqlx.DB) repository.Status {
	return &status{db: db}
}

// statusを投稿
func (r *status) Insert(ctx context.Context, status object.Status, mediaIDs []object.AttachmentID) (object.StatusID, error) {
	tx, err := r.db.BeginTxx(ctx, nil)
	if err != nil {
		return -1, fmt.Errorf("%w", err)
	}

	query := "INSERT INTO status (content, account_id) VALUES(?, ?)"
	row, err := tx.ExecContext(ctx, query, status.Content, status.Account.ID)
	if err != nil {
		return -1, fmt.Errorf("%w", err)
	}

	statusID, err := row.LastInsertId()
	if err != nil {
		return -1, fmt.Errorf("%w", err)
	}

	for _, mediaID := range mediaIDs {
		query = "INSERT INTO status_contain_attachment (status_id, attachment_id) VALUES(?, ?)"
		_, err := tx.ExecContext(ctx, query, statusID, mediaID)
		if err != nil {
			tx.Rollback()
			return -1, err
		}
	}
	err = tx.Commit()
	return statusID, err
}

// idからstatusを取得
func (r *status) FindByID(ctx context.Context, id object.StatusID) (*object.Status, error) {
	entity := new(object.Status)
	const query = `
SELECT
	s.id,
	s.content,
	s.create_at,
	a.id AS "account.id",
	a.username AS "account.username",
	a.password_hash AS "account.password_hash",
	a.create_at AS "account.create_at",
	CASE
		WHEN NOT EXISTS (
			SELECT
				*
			FROM
				relation AS r
				INNER JOIN account AS a ON r.following_id = a.id
				INNER JOIN status AS s ON s.account_id = a.id
			WHERE
				s.id = ?
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
						INNER JOIN status s ON s.account_id = a.id
					WHERE
						s.id = ?
				)
			GROUP BY
				following_id
		)
	END AS "account.followingcount",
	CASE
		WHEN NOT EXISTS (
			SELECT
				*
			FROM
				relation AS r
				INNER JOIN account AS a ON r.follower_id = a.id
				INNER JOIN status AS s ON s.account_id = a.id
			WHERE
				s.id = ?
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
					from
						account a
						INNER JOIN status s ON s.account_id = a.id
					WHERE
						s.id = ?
				)
			GROUP BY
				following_id
		)
	END AS "account.followerscount"
FROM
	status AS s
	JOIN account AS a ON s.account_id = a.id
WHERE
	s.id = ?
	`

	err := r.db.QueryRowxContext(ctx, query, id, id, id, id, id).StructScan(entity)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("%w", err)
	}
	return entity, nil
}

// idで指定したstatusを削除
func (r *status) Delete(ctx context.Context, id object.StatusID) error {
	const query = "DELETE FROM status WHERE id = ?"

	_, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil
		}
		return fmt.Errorf("%w", err)
	}
	return nil
}

// public timelineを取得
func (r *status) PublicTimeline(ctx context.Context, p object.Parameters) (object.Timelines, error) {
	var public object.Timelines
	var onlyMedia string
	if p.OnlyMedia {
		onlyMedia = "AND EXISTS(SELECT * FROM status_contain_attachment sca WHERE sca.status_id = s.id)"
	}

	query := fmt.Sprintf(`
SELECT
	s.id AS 'id',
	s.account_id AS 'account.id',
	s.create_at AS 'create_at',
	s.content AS 'content',
	a.username AS 'account.username',
	a.create_at AS 'account.create_at',
	CASE
		WHEN NOT EXISTS (
			SELECT
				*
			FROM
				relation AS r
				INNER JOIN account AS a ON r.following_id = a.id
				INNER JOIN status ON status.account_id = a.id
			WHERE
				status.id = s.id
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
						INNER JOIN status ON status.account_id = a.id
					WHERE
						status.id = s.id
				)
			GROUP BY
				following_id
		)
	END AS "account.followingcount",
	CASE
		WHEN NOT EXISTS (
			SELECT
				*
			FROM
				relation AS r
				INNER JOIN account AS a ON r.follower_id = a.id
				INNER JOIN status ON status.account_id = a.id
			WHERE
				status.id = s.id
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
						INNER JOIN status ON status.account_id = a.id
					WHERE
						status.id = s.id
				)
			GROUP BY
				follower_id
		)
	END AS "account.followerscount"
FROM
	status AS s
	JOIN account AS a ON s.account_id = a.id
WHERE
	s.id < ?
	AND s.id > ?
	%s
ORDER BY
	s.id
LIMIT
	?
	`, onlyMedia)

	err := r.db.SelectContext(ctx, &public, query, p.MaxID, p.SinceID, p.Limit)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("%w", err)
	}

	return public, nil
}

// home timelineを取得
func (r *status) HomeTimeline(ctx context.Context, loginID object.AccountID, p object.Parameters) (object.Timelines, error) {
	var home object.Timelines
	var onlyMedia string
	if p.OnlyMedia {
		onlyMedia = "AND EXISTS(SELECT * FROM status_contain_attachment sca WHERE sca.status_id = s.id)"
	}

	query := fmt.Sprintf(`
SELECT
	s.id,
	s.content,
	s.create_at,
	a.id AS "account.id",
	a.username AS "account.username",
	a.create_at AS "account.create_at",
	CASE
		WHEN NOT EXISTS (
			SELECT
				*
			FROM
				relation AS r
				INNER JOIN account AS a ON r.following_id = a.id
				INNER JOIN status ON status.account_id = a.id
			WHERE
				status.id = s.id
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
						INNER JOIN status ON status.account_id = a.id
					WHERE
						status.id = s.id
				)
			GROUP BY
				following_id
		)
	END AS "account.followingcount",
	CASE
		WHEN NOT EXISTS (
			SELECT
				*
			FROM
				relation AS r
				INNER JOIN account AS a ON r.follower_id = a.id
				INNER JOIN status ON status.account_id = a.id
			WHERE
				status.id = s.id
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
						INNER JOIN status ON status.account_id = a.id
					WHERE
						status.id = s.id
				)
			GROUP BY
				follower_id
		)
	END AS "account.followerscount"
FROM
	status AS s
	INNER JOIN (
		SELECT
			account.id,
			account.username,
			account.display_name,
			account.header,
			account.note,
			account.create_at
		FROM
			account
			INNER JOIN relation ON account.id = relation.follower_id
		WHERE
			account.id = ?
			OR relation.following_id = ?
	) AS a ON a.id = s.account_id
WHERE
	s.id > ?
	AND s.id < ?
	%s
ORDER BY
	s.id
LIMIT
	?
	`, onlyMedia)

	err := r.db.SelectContext(ctx, &home, query, loginID, loginID, p.SinceID, p.MaxID, p.Limit)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("%w", err)
	}

	return home, nil
}
