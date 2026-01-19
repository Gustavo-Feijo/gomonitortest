package user_test

import (
	"gomonitor/internal/testutil"
	"log"
	"os"
	"testing"

	"gorm.io/gorm"
)

var testDB *gorm.DB

func TestMain(m *testing.M) {
	db, cleanup, err := testutil.NewTestDBConnection()
	if err != nil {
		log.Fatal(err)
	}
	testDB = db
	code := m.Run()
	cleanup()
	os.Exit(code)
}

func setupTx(t *testing.T) *gorm.DB {
	t.Helper()
	tx := testDB.Begin()
	t.Cleanup(func() { tx.Rollback() })
	return tx
}
