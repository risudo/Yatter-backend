package dao

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"
	"yatter-backend-go/app/domain/object"
	"yatter-backend-go/app/domain/repository"

	"github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
)

const mysqlForeignKeyErrNo = 1452

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
//TODO: attachmentのidがなかった場合にロールバックしてstatusbadrequestにする
func (r *status) Insert(ctx context.Context, status object.Status, mediaIDs []object.AttachmentID) (object.StatusID, error) {
	query := "INSERT INTO status (content, account_id) VALUES(?, ?)"

	row, err := r.db.ExecContext(ctx, query, status.Content, status.Account.ID)
	if err != nil {
		return -1, fmt.Errorf("%w", err)
	}

	statusID, err := row.LastInsertId()
	if err != nil {
		return -1, fmt.Errorf("%w", err)
	}

	for _, mediaID := range mediaIDs {
		query = "INSERT INTO status_contain_attachment (status_id, attachment_id) VALUES(?, ?)"
		_, err := r.db.ExecContext(ctx, query, statusID, mediaID)
		if err != nil {
			if mysqlErr, ok := err.(*mysql.MySQLError); ok && mysqlErr.Number == mysqlForeignKeyErrNo {
				log.Println("👺media ids not found")
				return -1, err
			}
			return -1, err
		}
	}
	return statusID, nil
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
		a.create_at AS "account.create_at"
	FROM
		status AS s
	JOIN account AS a ON s.account_id = a.id
	WHERE s.id = ?`

	err := r.db.QueryRowxContext(ctx, query, id).StructScan(entity)
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
	const query = `
	SELECT
		s.id AS 'id',
		s.account_id AS 'account.id',
		s.create_at AS 'create_at',
		s.content AS 'content',
		a.username AS 'account.username',
		a.create_at AS 'account.create_at'
	FROM
		status AS s
	JOIN account AS a ON s.account_id = a.id
	WHERE s.id < ? AND s.id > ?
	ORDER BY s.id
	LIMIT ?;`

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
	const query = `
	SELECT
		s.id AS 'id',
		s.account_id AS 'account.id',
		a.username AS 'account.username',
		s.create_at AS 'create_at',
		s.content AS 'content',
		a.create_at AS 'account.create_at'
	FROM
		status AS s
	JOIN account AS a
	ON s.account_id = a.id
	JOIN relation
	ON a.id = relation.follower_id
	WHERE
		a.id = ?
	OR a.id
		IN (SELECT relation.follower_id
				FROM relation
				WHERE relation.following_id = ?)
	AND s.id < ? AND s.id > ?
	ORDER BY s.id
	LIMIT ?;`

	err := r.db.SelectContext(ctx, &home, query, loginID, loginID, p.MaxID, p.SinceID, p.Limit)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("%w", err)
	}

	return home, nil
}
