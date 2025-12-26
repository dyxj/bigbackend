//go:build integration

package integration

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/dyxj/bigbackend/internal/app"
	"github.com/dyxj/bigbackend/internal/config"
	"github.com/dyxj/bigbackend/internal/userprofile"
	"github.com/dyxj/bigbackend/pkg/logx"
	"github.com/dyxj/bigbackend/pkg/testx"
	"github.com/dyxj/bigbackend/test/faker"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestUserProfileHandler(t *testing.T) {
	cfg, err := config.LoadConfig()
	if err != nil {
		t.Fatalf("failed to load config: %v", err)
	}

	logger, err := logx.InitLogger()
	if err != nil {
		t.Fatalf("failed to init logger: %v", err)
	}

	dbConn := testx.GlobalEnv().DBConn()

	t.Cleanup(func() {
		truncateUserProfile(dbConn)
	})

	// TODO move server setup to test main
	srv := app.NewServer(logger, dbConn, cfg.HTTPServerConfig)

	testSrv := httptest.NewServer(srv.BuildRouter())
	defer testSrv.Close()

	payload := faker.UserProfileCreateRequest()

	var buf bytes.Buffer
	err = json.NewEncoder(&buf).Encode(&payload)
	if err != nil {
		t.Fatalf("failed to encode payload: %v", err)
	}

	request, err := http.NewRequest(
		"POST",
		buildUserProfileUrl(testSrv.URL, payload.UserID),
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
	var result userprofile.Response
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

func buildUserProfileUrl(url string, userId uuid.UUID) string {
	return fmt.Sprintf("%s/user/%s/profile", url, userId.String())
}
