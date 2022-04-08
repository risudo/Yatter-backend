package dao_test

import (
	"context"
	"fmt"
	"log"
	"testing"
	"yatter-backend-go/app/config"
	"yatter-backend-go/app/dao"
	"yatter-backend-go/app/domain/object"
	"yatter-backend-go/app/domain/repository"

	"github.com/jmoiron/sqlx"
)

// 実装案
// 1. トランザクションで頑張る
// 2. 別のdbを用意する

const notExistingUser = "notexist"

type mockdao struct {
	db *sqlx.DB
}

func (m *mockdao) Account() repository.Account {
	return dao.NewAccount(m.db)
}

func initMockDB(config dao.DBConfig) (*sqlx.DB, error) {
	driverName := "mysql"
	db, err := sqlx.Open(driverName, config.FormatDSN())
	if err != nil {
		return nil, fmt.Errorf("sqlx.Open failed: %w", err)
	}

	return db, nil
}

func setupDB() (*sqlx.DB, *sqlx.Tx, error) {
	daoCfg := config.MySQLConfig()
	db, err := initMockDB(daoCfg)
	if err != nil {
		return nil, nil, err
	}
	// トランザクション開始
	tx, _ := db.Beginx()
	// テーブルリセット
	if _, err := db.Exec("SET FOREIGN_KEY_CHECKS=0"); err != nil {
		return nil, nil, err
	}
	for _, table := range []string{"account", "status", "relation", "attachment"} {
		log.Println("table:", table)
		if _, err := db.Exec("DELETE FROM " + table); err != nil {
			return nil, nil, err
		}
	}
	return db, tx, nil
}

func TestAccount(t *testing.T) {
	db, tx, err := setupDB()
	defer tx.Rollback()

	m := &mockdao{db: db}
	repo := m.Account()
	ctx := context.Background()

	ra, err := repo.FindByUsername(ctx, notExistingUser)
	if err != nil || ra != nil {
		t.Fatal(err)
	}

	account := &object.Account{
		Username: "Michael",
	}
	err = repo.Insert(ctx, *account)
	if err != nil {
		t.Fatal(err)
	}
	ra, err = repo.FindByUsername(ctx, account.Username)
	if ra.Username != account.Username {
		t.Fatal(fmt.Errorf("does not match username"))
	}

	displayName := "Make"
	note := "note"
	account.DisplayName = &displayName
	account.Note = &note

	err = repo.Update(ctx, *account)
	ra, err = repo.FindByUsername(ctx, account.Username)
	if *ra.DisplayName != *account.DisplayName {
		t.Fatal(fmt.Errorf("does not match displayName"))
	}
}

func TestFindByUsername(t *testing.T) {
}
