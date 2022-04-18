package dao_test

import (
	"context"
	"fmt"
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

var preparedAccount = &object.Account{
	Username: "Michael",
}

var preparedStatus = &object.Status{
	Account: preparedAccount,
	Content: "content",
}

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
		if _, err := db.Exec("DELETE FROM " + table); err != nil {
			return nil, nil, err
		}
	}
	mockdao := &mockdao{db: db}
	ctx := context.Background()
	preparedAccount.ID, err = mockdao.Account().Insert(ctx, *preparedAccount)
	if err != nil {
		tx.Rollback()
		return nil, nil, err
	}
	preparedStatus.ID, err = mockdao.Status().Insert(ctx, *preparedStatus, nil)
	return mockdao, tx, nil
}

func TestFindByUsername(t *testing.T) {
	m, tx, err := setupDB()
	if err != nil {
		t.Fatal(err)
	}
	defer tx.Rollback()

	repo := m.Account()
	ctx := context.Background()

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
			userName:      preparedAccount.Username,
			expectAccount: preparedAccount,
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
			opt := cmpopts.IgnoreFields(object.Account{}, "CreateAt")
			if d := cmp.Diff(actual, tt.expectAccount, opt); len(d) != 0 {
				t.Fatalf("differs: (-got +want)\n%s", d)
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

	displayName := "Mike"
	note := "note"

	tests := []struct {
		name          string
		account       *object.Account
		expectErr     bool
		expectAccount *object.Account
	}{
		{
			name: "Update",
			account: &object.Account{
				ID:          preparedAccount.ID,
				Username:    preparedAccount.Username,
				DisplayName: &displayName,
				Note:        &note,
			},
			expectErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var gotErr bool
			err := repo.Update(ctx, *tt.account)

			if err != nil {
				gotErr = true
			}
			if gotErr != tt.expectErr {
				if err != nil {
					t.Fatal(err)
				} else {
					t.Fatal("expect error")
				}
			}

			updated, err := repo.FindByUsername(ctx, preparedAccount.Username)
			opt := cmpopts.IgnoreFields(object.Account{}, "CreateAt")
			if d := cmp.Diff(updated, tt.account, opt); len(d) != 0 {
				tx.Rollback()
				t.Fatalf("differs: (-got +want)\n%s", d)
			}
		})
	}
}

func TestStatusInsert(t *testing.T) {
	//note: 存在しないmediaIDを渡した時の処理が確定してないので保留
}

func TestStatusFindByID(t *testing.T) {
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
			id:           -100,
			expectStatus: nil,
		},
		{
			name:         "FindPreparedStatus",
			id:           preparedStatus.ID,
			expectStatus: preparedStatus,
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
			opt := cmpopts.IgnoreFields(object.Status{}, "CreateAt", "Account")
			if d := cmp.Diff(actual, tt.expectStatus, opt); len(d) != 0 {
				t.Fatalf("differs: (-got +want)\n%s", d)
			}
		})
	}
}

func TestStatusDelete(t *testing.T) {
	m, tx, err := setupDB()
	if err != nil {
		t.Fatal(err)
	}
	defer tx.Rollback()

	repo := m.Status()
	ctx := context.Background()

	tests := []struct {
		name        string
		id          object.StatusID
		expectIsErr bool
	}{
		{
			name:        "delete",
			id:          preparedStatus.ID,
			expectIsErr: false,
		},
		{
			name:        "deleteNotExist",
			id:          preparedStatus.ID,
			expectIsErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := repo.Delete(ctx, tt.id)
			if tt.expectIsErr == false && err != nil {
				t.Fatal(err)
			}
		})
	}
}

func TestStatusPublicTimeline(t *testing.T) {

}
