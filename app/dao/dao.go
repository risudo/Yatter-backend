package dao

import (
	"fmt"
	"log"
	"time"
	"yatter-backend-go/app/domain/repository"

	"github.com/jmoiron/sqlx"
)

type (
	// DAO interface
	Dao interface {
		// Get account repository
		Account() repository.Account

		// Get status repository
		Status() repository.Status

		// Get relation repository
		Relation() repository.Relation

		// Get attachment repository
		Attachment() repository.Attachment

		// Clear all data in DB
		InitAll() error
	}

	// Implementation for DAO
	dao struct {
		db *sqlx.DB
	}
)

// Create DAO
func New(config DBConfig) (Dao, error) {
	db, err := initDb(config)
	if err != nil {
		return nil, err
	}
	const connections = 10
	db.SetMaxIdleConns(connections)
	db.SetMaxOpenConns(connections)
	db.SetConnMaxLifetime(connections * time.Second)

	return &dao{db: db}, nil
}

func (d *dao) Account() repository.Account {
	return NewAccount(d.db)
}

func (d *dao) Status() repository.Status {
	return NewStatus(d.db)
}

func (d *dao) Relation() repository.Relation {
	return NewRelation(d.db)
}

func (d *dao) Attachment() repository.Attachment {
	return NewAttachment(d.db)
}

func (d *dao) InitAll() error {
	if err := d.exec("SET FOREIGN_KEY_CHECKS=0"); err != nil {
		return fmt.Errorf("can't disable FOREIGN_KEY_CHECKS: %w", err)
	}

	defer func() {
		err := d.exec("SET FOREIGN_KEY_CHECKS=0")
		if err != nil {
			log.Printf("Can't restore FOREIGN_KEY_CHECKS: %+v", err)
		}
	}()

	for _, table := range []string{"account", "status", "relation", "attachment"} {
		if err := d.exec("TRUNCATE TABLE " + table); err != nil {
			return fmt.Errorf("Can't truncate table "+table+": %w", err)
		}
	}

	return nil
}

func (d *dao) exec(query string, args ...interface{}) error {
	_, err := d.db.Exec(query, args...)
	return err
}
