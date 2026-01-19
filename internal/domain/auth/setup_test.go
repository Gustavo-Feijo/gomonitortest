package auth_test

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
	defer cleanup()
	testDB = db
	code := m.Run()
	cleanup()
	os.Exit(code)
}
