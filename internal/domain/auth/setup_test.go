package auth_test

import (
	"gomonitor/internal/testutil"
	"os"
	"testing"

	"gorm.io/gorm"
)

var testDB *gorm.DB

func TestMain(m *testing.M) {
	db, cleanup := testutil.NewTestConnection()
	testDB = db
	code := m.Run()
	cleanup()
	os.Exit(code)
}
