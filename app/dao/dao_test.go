package dao

import (
	"context"
	"testing"
	"yatter-backend-go/app/config"
	"yatter-backend-go/app/domain/object"
)

func setup() Dao {
	daoCfg := config.MySQLConfig()
	dao, err := New(daoCfg)
	if err != nil {
		panic(err)
	}

	err = dao.InitAll()
	if err != nil {
		panic(err)
	}

	return dao
}

// アカウントが正常に作成できているか
// create atがそれっぽい値になっているか
func TestAccountCreate(t *testing.T) {
	dao := setup()

	a := dao.Account()
	account := &object.Account{
		Username:     "testuser",
		PasswordHash: "testpass",
	}

	err := a.Create(context.Background(), account)
	if err != nil {
		t.Fatal(err)
	}
	if account == nil {
		t.Fatal(err)
	}
}

// FindByUsername
// userないときにnilが返ってくるか
// userいるときにentityが返ってくるか
func TestFindByUsername(t *testing.T) {
	dao := setup()
	ctx := context.Background()

	a := dao.Account()
	account, err := a.FindByUsername(ctx, "nosuchusername")
	if err != nil {
		t.Fatal(err)
	}
	if account != nil {
		t.Fatal(err)
	}
}

// status
// statusを投稿できるか
// create atがそれっぽい値になっているか
func TestStatusPost(t *testing.T) {
	dao := setup()

	_ = dao.Status()
}

// findbyid
// statusないときにnil返ってくるか
// statusあるときにentity返ってくるか

// delete
