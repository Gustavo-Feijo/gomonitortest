package main

import (
	"context"
	"gomonitor/internal/testutil"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

// Run the app, only starting the database.
// The app should be able to run normally, since the only dependency there is the database.
func TestRun(t *testing.T) {
	testutil.StartTestDB(t)

	ctx, cancel := context.WithCancel(t.Context())

	// Cancel the context after sometime to force the server to stop.
	go func() {
		time.Sleep(500 * time.Millisecond)
		cancel()
	}()

	runErr := Run(ctx)
	assert.Nil(t, runErr)
}

func TestRunNoDB(t *testing.T) {
	postgresContainer := testutil.StartTestDB(t)
	if err := postgresContainer.Terminate(t.Context()); err != nil {
		t.Logf("cleanup: failed to terminate postgres container: %v", err)
	}

	runErr := Run(t.Context())
	assert.NotNil(t, runErr)
}
