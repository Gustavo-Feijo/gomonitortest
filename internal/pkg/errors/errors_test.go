package errors_test

import (
	"errors"
	pkgerrors "gomonitor/internal/pkg/errors"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewBadRequestError(t *testing.T) {
	t.Parallel()
	underlying := errors.New("invalid input")

	err := pkgerrors.NewBadRequestError("bad request", underlying)

	if err == nil {
		t.Fatal("expected error, got nil")
	}

	if err.StatusCode != http.StatusBadRequest {
		t.Errorf("expected status %d, got %d", http.StatusBadRequest, err.StatusCode)
	}

	if err.Code != "BAD_REQUEST" {
		t.Errorf("expected code BAD_REQUEST, got %s", err.Code)
	}

	if err.Message != "bad request" {
		t.Errorf("expected message 'bad request', got %s", err.Message)
	}

	if err.Err != underlying {
		t.Errorf("expected underlying error to be preserved")
	}

	if err.File == "" || err.Line == 0 {
		t.Errorf("expected file and line to be set, got %s:%d", err.File, err.Line)
	}

	assert.Equal(t, "bad request", err.Error())
}

func TestNewUnauthorizedError(t *testing.T) {
	t.Parallel()
	err := pkgerrors.NewUnauthorizedError("unauthorized")

	if err.StatusCode != http.StatusUnauthorized {
		t.Errorf("expected status %d, got %d", http.StatusUnauthorized, err.StatusCode)
	}

	if err.Code != "UNAUTHORIZED" {
		t.Errorf("expected code UNAUTHORIZED, got %s", err.Code)
	}
}

func TestNewNotFoundError(t *testing.T) {
	t.Parallel()
	err := pkgerrors.NewNotFoundError("not found")

	if err.StatusCode != http.StatusNotFound {
		t.Errorf("expected status %d, got %d", http.StatusNotFound, err.StatusCode)
	}

	if err.Code != "NOT_FOUND" {
		t.Errorf("expected code NOT_FOUND, got %s", err.Code)
	}
}

func TestNewConflictError(t *testing.T) {
	t.Parallel()
	err := pkgerrors.NewConflictError("conflict")

	if err.StatusCode != http.StatusConflict {
		t.Errorf("expected status %d, got %d", http.StatusConflict, err.StatusCode)
	}

	if err.Code != "CONFLICT" {
		t.Errorf("expected code CONFLICT, got %s", err.Code)
	}
}

func TestNewInternalError_DefaultMessage(t *testing.T) {
	t.Parallel()
	underlying := errors.New("db down")

	err := pkgerrors.NewInternalError(underlying)

	if err.StatusCode != http.StatusInternalServerError {
		t.Errorf("expected status %d, got %d", http.StatusInternalServerError, err.StatusCode)
	}

	if err.Code != "INTERNAL_ERROR" {
		t.Errorf("expected code INTERNAL_ERROR, got %s", err.Code)
	}

	if err.Message != "An unexpected error occurred" {
		t.Errorf("unexpected message: %s", err.Message)
	}

	if err.Err != underlying {
		t.Errorf("expected underlying error to be preserved")
	}
}

func TestNewForbiddenError(t *testing.T) {
	t.Parallel()
	err := pkgerrors.NewForbiddenError()

	if err.StatusCode != http.StatusForbidden {
		t.Errorf("expected status %d, got %d", http.StatusForbidden, err.StatusCode)
	}

	if err.Code != "FORBIDDEN" {
		t.Errorf("expected code FORBIDDEN, got %s", err.Code)
	}

	if err.Message != "FORBIDDEN" {
		t.Errorf("expected message FORBIDDEN, got %s", err.Message)
	}
}

func TestNewTooManyRequestError(t *testing.T) {
	t.Parallel()
	err := pkgerrors.NewTooManyRequestsError("too many")

	if err.StatusCode != http.StatusTooManyRequests {
		t.Errorf("expected status %d, got %d", http.StatusTooManyRequests, err.StatusCode)
	}

	if err.Code != "TOO_MANY_REQUEST" {
		t.Errorf("expected code TOO_MANY_REQUEST, got %s", err.Code)
	}
}

// Just so i can get my sweet 100% coverage
func TestCallerFailureCoverage(t *testing.T) {
	t.Parallel()
	orig := pkgerrors.RuntimeCaller
	defer func() { pkgerrors.RuntimeCaller = orig }()

	pkgerrors.RuntimeCaller = func(skip int) (uintptr, string, int, bool) {
		return 0, "", 0, false
	}

	err := pkgerrors.NewBadRequestError("test")

	if err.File != "unknown" {
		t.Fatalf("expected unknown file")
	}
}
