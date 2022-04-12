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

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/jmoiron/sqlx"
)

//TODO:失敗したときもロールバックさせる

const notExistingUser = "notexist"

type mockdao struct {
	db *sqlx.DB
}

func (m *mockdao) Account() repository.Account {
	return dao.NewAccount(m.db)
}

func (m *mockdao) Status() repository.Status {
	return dao.NewStatus(m.db)
}

func initMockDB(config dao.DBConfig) (*sqlx.DB, error) {
	driverName := "mysql"
	db, err := sqlx.Open(driverName, config.FormatDSN())
	if err != nil {
		return nil, fmt.Errorf("sqlx.Open failed: %w", err)
	}

	return db, nil
}

func setupDB() (*mockdao, *sqlx.Tx, error) {
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
	for _, table := range []string{"account", "status", "relation", "attachment", "status_contain_attachment"} {
		log.Println("table:", table)
		if _, err := db.Exec("DELETE FROM " + table); err != nil {
			return nil, nil, err
		}
	}
	return &mockdao{db: db}, tx, nil
}

func TestFindByUsername(t *testing.T) {
	m, tx, err := setupDB()
	if err != nil {
		t.Fatal(err)
	}
	defer tx.Rollback()

	repo := m.Account()
	ctx := context.Background()

	account := &object.Account{
		Username: "Michael",
	}

	err = repo.Insert(ctx, *account)
	if err != nil {
		tx.Rollback()
		t.Fatal(err)
	}

	tests := []struct {
		name          string
		userName      string
		expectAccount *object.Account
	}{
		{
			name:          "NotExistingUser",
			userName:      notExistingUser,
			expectAccount: nil,
		},
		{
			name:          "ExistingUser",
			userName:      account.Username,
			expectAccount: account,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			actual, err := repo.FindByUsername(ctx, tt.userName)
			if err != nil {
				t.Fatal(err)
			}
			if actual == nil && actual == tt.expectAccount {
				return
			}
			opt := cmpopts.IgnoreFields(object.Account{}, "CreateAt", "ID")
			if d := cmp.Diff(actual, tt.expectAccount, opt); len(d) != 0 {
				t.Errorf("differs: (-got +want)\n%s", d)
			}
		})
	}
}

func TestAccountUpdate(t *testing.T) {
	m, tx, err := setupDB()
	if err != nil {
		t.Fatal(err)
	}
	defer tx.Rollback()

	repo := m.Account()
	ctx := context.Background()

	account := &object.Account{
		Username: "Michael",
	}
	err = repo.Insert(ctx, *account)
	if err != nil {
		tx.Rollback()
		t.Fatal(err)
	}

	displayName := "Mike"
	note := "note"
	account.DisplayName = &displayName
	account.Note = &note

	err = repo.Update(ctx, *account)
	updated, err := repo.FindByUsername(ctx, account.Username)
	opt := cmpopts.IgnoreFields(object.Account{}, "CreateAt", "ID")
	if d := cmp.Diff(updated, account, opt); len(d) != 0 {
		tx.Rollback()
		t.Errorf("differs: (-got +want)\n%s", d)
	}
}

func TestStatus(t *testing.T) {
	m, tx, err := setupDB()
	if err != nil {
		t.Fatal(err)
	}
	defer tx.Rollback()

	repo := m.Status()
	ctx := context.Background()

	tests := []struct {
		name         string
		id           object.StatusID
		expectStatus *object.Status
	}{
		{
			name:         "FindNotExistingID",
			id:           100,
			expectStatus: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			actual, err := repo.FindByID(ctx, tt.id)
			if err != nil {
				t.Fatal(err)
			}
			if actual == nil && actual == tt.expectStatus {
				return
			}
			opt := cmpopts.IgnoreFields(object.Status{}, "CreateAt")
			if d := cmp.Diff(actual, tt.expectStatus, opt); len(d) != 0 {
				t.Errorf("differs: (-got +want)\n%s", d)
			}
		})
	}

}
