//go:build integration

package integration

import (
	"bytes"
	"context"
	"encoding/json"
	"log"
	"net/http"
	"testing"
	"time"

	"cloud.google.com/go/civil"
	"github.com/dyxj/bigbackend/internal/sqlgen/bigbackend/public/entity"
	"github.com/dyxj/bigbackend/internal/user/profile"
	"github.com/dyxj/bigbackend/pkg/httpx"
	"github.com/dyxj/bigbackend/pkg/testx"
	"github.com/dyxj/bigbackend/test/faker"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestUserProfileCreatorHandler_ShouldCreate(t *testing.T) {
	testSrv := testx.GlobalEnv().HttpTestServer()

	dbConn := testx.GlobalEnv().DBConn()
	t.Cleanup(func() {
		truncateUserProfile(dbConn)
	})

	payload := faker.UserProfileCreateRequest()

	var buf bytes.Buffer
	err := json.NewEncoder(&buf).Encode(&payload)
	if err != nil {
		t.Fatalf("failed to encode payload: %v", err)
	}

	request, err := http.NewRequest(
		"POST",
		buildUserProfileUrl(testSrv.URL, payload.UserID.String()),
		&buf)
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
	var result profile.Response
	err = json.NewDecoder(resp.Body).
		Decode(&result)
	if err != nil {
		t.Fatalf("failed to decode response body: %v", err)
	}

	assert.Equal(t, http.StatusCreated, resp.StatusCode)

	assert.NotEqual(t, uuid.Nil, result.ID)
	assert.Equal(t, payload.UserID, result.UserID)
	assert.Equal(t, payload.FirstName, result.FirstName)
	assert.Equal(t, payload.LastName, result.LastName)
	assert.Equal(t, payload.DateOfBirth, result.DateOfBirth)
	assert.False(t, result.CreateTime.IsZero())
	assert.False(t, result.UpdateTime.IsZero())
	assert.Equal(t, int32(1), result.Version)
}

func TestUserProfileCreatorHandler_PayloadValidationError(t *testing.T) {
	testSrv := testx.GlobalEnv().HttpTestServer()

	dbConn := testx.GlobalEnv().DBConn()
	t.Cleanup(func() {
		truncateUserProfile(dbConn)
	})

	ttc := []struct {
		name    string
		mod     func(*profile.CreateRequest)
		errResp httpx.ErrorResponse
	}{
		{
			name: "missing user ID",
			mod: func(r *profile.CreateRequest) {
				r.UserID = uuid.Nil
			},
			errResp: httpx.ErrorResponse{
				Code:    httpx.CodeBadRequest,
				Message: "validation failed",
				Details: map[string]string{
					"userId": "is required",
				},
			},
		},
		{
			name: "missing first name",
			mod: func(r *profile.CreateRequest) {
				r.FirstName = ""
			},
			errResp: httpx.ErrorResponse{
				Code:    httpx.CodeBadRequest,
				Message: "validation failed",
				Details: map[string]string{
					"firstName": "is required",
				},
			},
		},
		{
			name: "missing last name",
			mod: func(r *profile.CreateRequest) {
				r.LastName = ""
			},
			errResp: httpx.ErrorResponse{
				Code:    httpx.CodeBadRequest,
				Message: "validation failed",
				Details: map[string]string{
					"lastName": "is required",
				},
			},
		},
		{
			name: "date of birth in the future",
			mod: func(input *profile.CreateRequest) {
				input.DateOfBirth = civil.DateOf(time.Now().Add(time.Hour * 24))
			},
			errResp: httpx.ErrorResponse{
				Code:    httpx.CodeBadRequest,
				Message: "validation failed",
				Details: map[string]string{
					"dateOfBirth": "is invalid or in the future",
				},
			},
		},
	}

	for _, tc := range ttc {
		t.Run(tc.name, func(t *testing.T) {
			payload := faker.UserProfileCreateRequest()
			tc.mod(&payload)

			var buf bytes.Buffer
			err := json.NewEncoder(&buf).Encode(&payload)
			if err != nil {
				t.Fatalf("failed to encode payload: %v", err)
			}

			request, err := http.NewRequest(
				"POST",
				buildUserProfileUrl(testSrv.URL, payload.UserID.String()),
				&buf)
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
			err = json.NewDecoder(resp.Body).Decode(&result)
			if err != nil {
				t.Fatalf("failed to decode response body: %v", err)
			}

			assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
			assert.Equal(t, tc.errResp, result)
		})
	}
}

func TestUserProfileCreatorHandler_InvalidJsonError(t *testing.T) {
	testSrv := testx.GlobalEnv().HttpTestServer()

	dbConn := testx.GlobalEnv().DBConn()
	t.Cleanup(func() {
		truncateUserProfile(dbConn)
	})

	ttc := []struct {
		name              string
		mod               func(*profile.CreateRequest)
		resultDetailError string
	}{
		{
			name: "zero date of birth",
			mod: func(input *profile.CreateRequest) {
				input.DateOfBirth = civil.Date{}
			},
			resultDetailError: "parsing time \"0000-00-00\": month out of range",
		},
		{
			name: "invalid date of birth",
			mod: func(input *profile.CreateRequest) {
				input.DateOfBirth = civil.Date{Year: 2024, Month: 13, Day: 32}
			},
			resultDetailError: "parsing time \"2024-13-32\": month out of range",
		},
	}

	for _, tc := range ttc {
		t.Run(tc.name, func(t *testing.T) {
			payload := faker.UserProfileCreateRequest()
			tc.mod(&payload)

			var buf bytes.Buffer
			err := json.NewEncoder(&buf).Encode(&payload)
			if err != nil {
				t.Fatalf("failed to encode payload: %v", err)
			}

			request, err := http.NewRequest(
				"POST",
				buildUserProfileUrl(testSrv.URL, payload.UserID.String()),
				&buf)
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
			err = json.NewDecoder(resp.Body).Decode(&result)
			if err != nil {
				t.Fatalf("failed to decode response body: %v", err)
			}

			assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
			assert.Equal(t, "invalid request body", result.Message)
			assert.Equal(t, result.Details["error"], tc.resultDetailError)
		})
	}
}

func TestUserProfileCreatorHandler_URL_ID_payload_ID_mismatch(t *testing.T) {
	testSrv := testx.GlobalEnv().HttpTestServer()

	dbConn := testx.GlobalEnv().DBConn()
	t.Cleanup(func() {
		truncateUserProfile(dbConn)
	})

	payload := faker.UserProfileCreateRequest()
	userIdMismatch := uuid.New()

	var buf bytes.Buffer
	err := json.NewEncoder(&buf).Encode(&payload)
	if err != nil {
		t.Fatalf("failed to encode payload: %v", err)
	}

	request, err := http.NewRequest(
		"POST",
		buildUserProfileUrl(testSrv.URL, userIdMismatch.String()),
		&buf)
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
	err = json.NewDecoder(resp.Body).Decode(&result)
	if err != nil {
		t.Fatalf("failed to decode response body: %v", err)
	}

	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	assert.Equal(t, "user ID in URL does not match user ID in request body", result.Message)
	assert.Nil(t, result.Details)
}

func TestUserProfileCreatorHandler_Unique_Violation(t *testing.T) {
	testSrv := testx.GlobalEnv().HttpTestServer()

	dbConn := testx.GlobalEnv().DBConn()
	t.Cleanup(func() {
		truncateUserProfile(dbConn)
	})

	payload := faker.UserProfileCreateRequest()

	existingUserProfile := entity.UserProfile{
		UserID:      payload.UserID,
		FirstName:   payload.FirstName,
		LastName:    payload.LastName,
		DateOfBirth: payload.DateOfBirth,
	}
	existingUserProfile, err := userProfileCreatorRepo.
		InsertUserProfile(context.Background(), dbConn, existingUserProfile)
	if err != nil {
		t.Fatalf("failed to insert existing user profile: %v", err)
	}

	var buf bytes.Buffer
	err = json.NewEncoder(&buf).Encode(&payload)
	if err != nil {
		t.Fatalf("failed to encode payload: %v", err)
	}

	request, err := http.NewRequest(
		"POST",
		buildUserProfileUrl(testSrv.URL, payload.UserID.String()),
		&buf)
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
	err = json.NewDecoder(resp.Body).Decode(&result)
	if err != nil {
		t.Fatalf("failed to decode response body: %v", err)
	}

	assert.Equal(t, http.StatusConflict, resp.StatusCode)
	assert.Equal(t, "user profile already exists", result.Message)
	assert.Nil(t, result.Details)
}
