package profile_test

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/dyxj/bigbackend/internal/user/profile"
	"github.com/dyxj/bigbackend/pkg/httpx"
	"github.com/dyxj/bigbackend/pkg/logx"
	"github.com/dyxj/bigbackend/test/faker"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestGetterHandler_InternalServerError(t *testing.T) {
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

	mapper := new(profile.UserProfileMapper)
	getterMock := new(faker.UserProfileGetterMock)
	getterHandler := profile.NewGetterHandler(logger, getterMock, mapper)

	userId := uuid.New()

	request := httptest.NewRequest(
		"GET",
		"/user/{id}/profile",
		nil,
	)
	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("id", userId.String())
	request = request.WithContext(
		context.WithValue(request.Context(), chi.RouteCtxKey, rctx),
	)

	rr := httptest.NewRecorder()

	getterMock.On("GetUserProfileByUserID", mock.Anything, userId).
		Return(profile.UserProfile{}, errors.New("fake error"))

	getterHandler.ServeHTTP(rr, request)

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
	getterMock.AssertNumberOfCalls(t, "GetUserProfileByUserID", 1)
}
