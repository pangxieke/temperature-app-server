package test

import (
	"testing"

	"github.com/jinzhu/gorm"
	_ "github.com/mattn/go-sqlite3"
	"gopkg.in/testfixtures.v2"
)

type fixtures struct {
	*testfixtures.Context
	*gorm.DB
}

func newFixtures() (*fixtures, error) {
	var err error
	db, err := gorm.Open("sqlite3", "file::memory:?cache=shared")
	if err != nil {
		return nil, err
	}

	testfixtures.SkipDatabaseNameCheck(true)
	f, err := testfixtures.NewFolder(db.DB(), &testfixtures.SQLite{}, "../test/fixtures")
	if err != nil {
		return nil, err
	}
	return &fixtures{f, db}, nil
}

func (f *fixtures) Prepare(t *testing.T) {
	//if err := f.Load(); err != nil {
	//	log.Fatalf("cannot load fixtures, err: %+v", err)
	//}
}
