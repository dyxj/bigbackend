//go:build integration

package integration

import (
	"encoding/json"
	"log"
	"net/http"
	"testing"

	"github.com/dyxj/bigbackend/internal/userprofile"
	"github.com/dyxj/bigbackend/pkg/httpx"
	"github.com/dyxj/bigbackend/pkg/testx"
	"github.com/dyxj/bigbackend/test/faker"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestUserProfileGetterHandler_Got(t *testing.T) {
	testSrv := testx.GlobalEnv().HttpTestServer()

	dbConn := testx.GlobalEnv().DBConn()
	t.Cleanup(func() {
		truncateUserProfile(dbConn)
	})

	ctx := t.Context()

	creator := userprofile.NewCreatorSQLDB(logger)
	uProfile := faker.UserProfileEntity()
	inserted, err := creator.InsertUserProfile(ctx, dbConn, uProfile)
	if err != nil {
		t.Fatalf("failed to insert user profile: %v", err)
	}

	request, err := http.NewRequest(
		"GET",
		buildUserProfileUrl(testSrv.URL, inserted.UserID.String()),
		nil)
	if err != nil {
		t.Fatalf("failed to build request: %v", err)
	}

	resp, err := testSrv.Client().Do(request)
	if err != nil {
		t.Fatalf("failed to perform request: %v", err)
	}
	defer func() {
		err := resp.Body.Close()
		if err != nil {
			log.Printf("failed to close response body: %v", err)
		}
	}()
	var result userprofile.Response
	err = json.NewDecoder(resp.Body).
		Decode(&result)
	if err != nil {
		t.Fatalf("failed to decode response body: %v", err)
	}

	assert.Equal(t, http.StatusOK, resp.StatusCode)

	assert.NotEqual(t, uuid.Nil, result.ID)
	assert.Equal(t, inserted.UserID, result.UserID)
	assert.Equal(t, inserted.FirstName, result.FirstName)
	assert.Equal(t, inserted.LastName, result.LastName)
	assert.Equal(t, inserted.DateOfBirth, result.DateOfBirth)
	assert.False(t, result.CreateTime.IsZero())
	assert.False(t, result.UpdateTime.IsZero())
	assert.Equal(t, int32(1), result.Version)
}

func TestUserProfileGetterHandler_NotFound(t *testing.T) {
	testSrv := testx.GlobalEnv().HttpTestServer()

	dbConn := testx.GlobalEnv().DBConn()
	t.Cleanup(func() {
		truncateUserProfile(dbConn)
	})

	notFoundUserId := uuid.New()

	request, err := http.NewRequest(
		"GET",
		buildUserProfileUrl(testSrv.URL, notFoundUserId.String()),
		nil)
	if err != nil {
		t.Fatalf("failed to build request: %v", err)
	}

	resp, err := testSrv.Client().Do(request)
	if err != nil {
		t.Fatalf("failed to perform request: %v", err)
	}
	defer func() {
		err := resp.Body.Close()
		if err != nil {
			log.Printf("failed to close response body: %v", err)
		}
	}()
	var result httpx.ErrorResponse
	err = json.NewDecoder(resp.Body).
		Decode(&result)
	if err != nil {
		t.Fatalf("failed to decode response body: %v", err)
	}

	assert.Equal(t, http.StatusNotFound, resp.StatusCode)

	assert.Equal(t, httpx.CodeEntityNotFound, result.Code)
	assert.Equal(t, "entity not found", result.Message)
	assert.Nil(t, result.Details)
}

func TestUserProfileGetterHandler_InvalidUserId(t *testing.T) {
	testSrv := testx.GlobalEnv().HttpTestServer()

	dbConn := testx.GlobalEnv().DBConn()
	t.Cleanup(func() {
		truncateUserProfile(dbConn)
	})

	request, err := http.NewRequest(
		"GET",
		buildUserProfileUrl(testSrv.URL, "invalid-uuid"),
		nil)
	if err != nil {
		t.Fatalf("failed to build request: %v", err)
	}

	resp, err := testSrv.Client().Do(request)
	if err != nil {
		t.Fatalf("failed to perform request: %v", err)
	}
	defer func() {
		err := resp.Body.Close()
		if err != nil {
			log.Printf("failed to close response body: %v", err)
		}
	}()
	var result httpx.ErrorResponse
	err = json.NewDecoder(resp.Body).
		Decode(&result)
	if err != nil {
		t.Fatalf("failed to decode response body: %v", err)
	}

	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)

	assert.Equal(t, httpx.CodeBadRequest, result.Code)
	assert.Equal(t, "invalid id", result.Message)
	assert.Equal(t, "invalid UUID length: 12", result.Details["error"])
}
