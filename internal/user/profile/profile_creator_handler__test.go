package profile_test

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"errors"
	"io"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/dyxj/bigbackend/internal/user/profile"
	"github.com/dyxj/bigbackend/pkg/errorx"
	"github.com/dyxj/bigbackend/pkg/httpx"
	"github.com/dyxj/bigbackend/pkg/logx"
	"github.com/dyxj/bigbackend/test/faker"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestCreatorHandler_BeginTxError(t *testing.T) {
	logger, err := logx.InitLogger()
	if err != nil {
		t.Fatalf("failed to initialize logger: %v", err)
	}

	dbMock, err := faker.NewTransactionManagerMock()
	if err != nil {
		t.Fatalf("failed to create transaction manager mock: %v", err)
	}
	defer func(dbMock *faker.TransactionManagerMock) {
		err := dbMock.Close()
		if err != nil {
			log.Printf("failed to close db mock: %v", err)
		}
	}(dbMock)

	creatorMock := new(faker.UserProfileCreatorMock)
	mapper := new(profile.UserProfileMapper)

	handler := profile.NewCreatorHandler(
		logger,
		dbMock,
		creatorMock,
		mapper,
	)

	payload := faker.UserProfileCreateRequest()
	var buf bytes.Buffer
	err = json.NewEncoder(&buf).Encode(&payload)
	if err != nil {
		t.Fatalf("failed to encode payload: %v", err)
	}

	request := httptest.NewRequest(
		"POST",
		"/user/{id}/profile",
		&buf,
	)
	request.SetPathValue("id", payload.UserID.String())
	rr := httptest.NewRecorder()

	dbMock.On("BeginTx", mock.Anything, mock.Anything).
		Return(&sql.Tx{}, errors.New("fake db error")).
		Once()

	handler.ServeHTTP(rr, request)

	expectedResultPayload := httpx.ErrorResponse{
		Code:    httpx.CodeServerError,
		Message: "internal server error",
		Details: nil,
	}

	result := rr.Result()

	var resultPayload httpx.ErrorResponse
	err = json.NewDecoder(result.Body).Decode(&resultPayload)
	if err != nil {
		t.Fatalf("failed to decode response body: %v", err)
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			log.Printf("failed to close response body: %v", err)
		}
	}(result.Body)

	assert.Equal(t, http.StatusInternalServerError, result.StatusCode)
	assert.Equal(t, expectedResultPayload, resultPayload)
	dbMock.AssertNumberOfCalls(t, "BeginTx", 1)
	creatorMock.AssertNotCalled(t, "CreateUserProfileTx")
}

func TestCreatorHandler_CreatorValidationError(t *testing.T) {
	logger, err := logx.InitLogger()
	if err != nil {
		t.Fatalf("failed to initialize logger: %v", err)
	}

	dbMock, err := faker.NewTransactionManagerMock()
	if err != nil {
		t.Fatalf("failed to create transaction manager mock: %v", err)
	}
	defer func(dbMock *faker.TransactionManagerMock) {
		err := dbMock.Close()
		if err != nil {

		}
	}(dbMock)

	creatorMock := new(faker.UserProfileCreatorMock)
	mapper := new(profile.UserProfileMapper)

	handler := profile.NewCreatorHandler(
		logger,
		dbMock,
		creatorMock,
		mapper,
	)

	payload := faker.UserProfileCreateRequest()
	var buf bytes.Buffer
	err = json.NewEncoder(&buf).Encode(&payload)
	if err != nil {
		t.Fatalf("failed to encode payload: %v", err)
	}

	request := httptest.NewRequest(
		"POST",
		"/user/{id}/profile",
		&buf,
	)
	request.SetPathValue("id", payload.UserID.String())
	rr := httptest.NewRecorder()

	dbMock.SqlMock().ExpectBegin()
	dbMock.On("BeginTx", mock.Anything, mock.Anything).
		Run(dbMock.ReturnTx)

	creatorMock.On("CreateUserProfileTx", mock.Anything, mock.Anything, mock.Anything).
		Return(profile.UserProfile{}, errors.New("fake unexpected error"))

	handler.ServeHTTP(rr, request)

	expectedResultPayload := httpx.ErrorResponse{
		Code:    httpx.CodeServerError,
		Message: "internal server error",
		Details: nil,
	}

	result := rr.Result()

	var resultPayload httpx.ErrorResponse
	err = json.NewDecoder(result.Body).Decode(&resultPayload)
	if err != nil {
		t.Fatalf("failed to decode response body: %v", err)
	}
	defer result.Body.Close()

	assert.Equal(t, http.StatusInternalServerError, result.StatusCode)
	assert.Equal(t, expectedResultPayload, resultPayload)
	dbMock.AssertNumberOfCalls(t, "BeginTx", 1)
	assert.NoError(t, dbMock.SqlMock().ExpectationsWereMet())
	creatorMock.AssertNumberOfCalls(t, "CreateUserProfileTx", 1)
}

func TestCreatorHandler_UnexpectedCreatorError(t *testing.T) {
	logger, err := logx.InitLogger()
	if err != nil {
		t.Fatalf("failed to initialize logger: %v", err)
	}

	dbMock, err := faker.NewTransactionManagerMock()
	if err != nil {
		t.Fatalf("failed to create transaction manager mock: %v", err)
	}
	defer func(dbMock *faker.TransactionManagerMock) {
		err := dbMock.Close()
		if err != nil {

		}
	}(dbMock)

	creatorMock := new(faker.UserProfileCreatorMock)
	mapper := new(profile.UserProfileMapper)

	handler := profile.NewCreatorHandler(
		logger,
		dbMock,
		creatorMock,
		mapper,
	)

	payload := faker.UserProfileCreateRequest()
	var buf bytes.Buffer
	err = json.NewEncoder(&buf).Encode(&payload)
	if err != nil {
		t.Fatalf("failed to encode payload: %v", err)
	}

	request := httptest.NewRequest(
		"POST",
		"/user/{id}/profile",
		&buf,
	)
	request.SetPathValue("id", payload.UserID.String())
	rr := httptest.NewRecorder()

	dbMock.SqlMock().ExpectBegin()
	dbMock.On("BeginTx", mock.Anything, mock.Anything).
		Run(dbMock.ReturnTx)

	creatorMock.On("CreateUserProfileTx", mock.Anything, mock.Anything, mock.Anything).
		Return(profile.UserProfile{}, &errorx.ValidationError{
			Properties: map[string]string{"firstName": "fake validation error"},
		})

	handler.ServeHTTP(rr, request)

	expectedResultPayload := httpx.ErrorResponse{
		Code:    httpx.CodeBadRequest,
		Message: "validation failed",
		Details: map[string]string{"firstName": "fake validation error"},
	}

	result := rr.Result()

	var resultPayload httpx.ErrorResponse
	err = json.NewDecoder(result.Body).Decode(&resultPayload)
	if err != nil {
		t.Fatalf("failed to decode response body: %v", err)
	}
	defer result.Body.Close()

	assert.Equal(t, http.StatusBadRequest, result.StatusCode)
	assert.Equal(t, expectedResultPayload, resultPayload)
	dbMock.AssertNumberOfCalls(t, "BeginTx", 1)
	assert.NoError(t, dbMock.SqlMock().ExpectationsWereMet())
	creatorMock.AssertNumberOfCalls(t, "CreateUserProfileTx", 1)
}

func TestCreatorHandler_TxCommitError(t *testing.T) {
	logger, err := logx.InitLogger()
	if err != nil {
		t.Fatalf("failed to initialize logger: %v", err)
	}

	dbMock, err := faker.NewTransactionManagerMock()
	if err != nil {
		t.Fatalf("failed to create transaction manager mock: %v", err)
	}
	defer func(dbMock *faker.TransactionManagerMock) {
		err := dbMock.Close()
		if err != nil {
			log.Printf("failed to close db mock: %v", err)
		}
	}(dbMock)

	creatorMock := new(faker.UserProfileCreatorMock)
	mapper := new(profile.UserProfileMapper)

	handler := profile.NewCreatorHandler(
		logger,
		dbMock,
		creatorMock,
		mapper,
	)

	payload := faker.UserProfileCreateRequest()
	var buf bytes.Buffer
	err = json.NewEncoder(&buf).Encode(&payload)
	if err != nil {
		t.Fatalf("failed to encode payload: %v", err)
	}

	request := httptest.NewRequest(
		"POST",
		"/user/{id}/profile",
		&buf,
	)
	request.SetPathValue("id", payload.UserID.String())
	rr := httptest.NewRecorder()

	dbMock.SqlMock().ExpectBegin()
	dbMock.On("BeginTx", mock.Anything, mock.Anything).
		Run(dbMock.ReturnCommitedTx).
		Once()

	creatorMock.On("CreateUserProfileTx", mock.Anything, mock.Anything, mock.Anything).
		Run(func(args mock.Arguments) {
			input := args.Get(2).(profile.UserProfile)
			creatorMock.ExpectedCalls[0].ReturnArguments = mock.Arguments{input, nil}
		})

	handler.ServeHTTP(rr, request)

	expectedResultPayload := httpx.ErrorResponse{
		Code:    httpx.CodeServerError,
		Message: "internal server error",
		Details: nil,
	}

	result := rr.Result()

	var resultPayload httpx.ErrorResponse
	err = json.NewDecoder(result.Body).Decode(&resultPayload)
	if err != nil {
		t.Fatalf("failed to decode response body: %v", err)
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			log.Printf("failed to close response body: %v", err)
		}
	}(result.Body)

	assert.Equal(t, http.StatusInternalServerError, result.StatusCode)
	assert.Equal(t, expectedResultPayload, resultPayload)
	dbMock.AssertNumberOfCalls(t, "BeginTx", 1)
	assert.NoError(t, dbMock.SqlMock().ExpectationsWereMet())
	creatorMock.AssertNumberOfCalls(t, "CreateUserProfileTx", 1)
}
