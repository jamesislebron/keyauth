package mysql_test

import (
	"github.com/defineiot/keyauth/dao/verifycode"
	"github.com/defineiot/keyauth/dao/verifycode/mysql"
	"github.com/defineiot/keyauth/internal/conf/mock"
)

func newTestStore() verifycode.Store {
	conf := mock.NewConfig()
	db, err := conf.GetDBConn()
	if err != nil {
		panic(err)
	}

	store, err := mysql.NewVerifyCodeStore(db)
	if err != nil {
		panic(err)
	}

	return store
}

type verifyCodeSuit struct {
	store  verifycode.Store
	code   *verifycode.VerifyCode
	target string
}

func (s *verifyCodeSuit) TearDown() {
	s.store.Close()
}

func (s *verifyCodeSuit) SetUp() {
	s.target = "18108053819"

	s.code = &verifycode.VerifyCode{
		Purpose:    verifycode.Registry,
		SendMode:   verifycode.Mobile,
		SendTarget: s.target,
	}

	s.store = newTestStore()

}