package userprofile_test

import (
	"context"
	"database/sql"
	"testing"
	"time"

	"cloud.google.com/go/civil"
	"github.com/dyxj/bigbackend/internal/faker"
	"github.com/dyxj/bigbackend/internal/sqlgen/bigbackend/public/entity"
	"github.com/dyxj/bigbackend/internal/userprofile"
	"github.com/dyxj/bigbackend/pkg/errorx"
	"github.com/dyxj/bigbackend/pkg/logx"
	"github.com/stretchr/testify/mock"
	"gotest.tools/v3/assert"
)

// Test that
// - userProfile is sanitized before creation
// - InsertUserProfile is called
// - created userProfile is returned
//
// - does not mock InsertUserProfile exact behaviour
func TestCreator_CreateUserProfileTx_Successfully(t *testing.T) {
	logger, err := logx.InitLogger()
	if err != nil {
		t.Fatalf("failed to initialize logger: %v", err)
	}

	mockRepo := new(mockCreatorRepo)
	mockRepo.
		On("InsertUserProfile", mock.Anything, mock.Anything, mock.Anything).
		Run(func(args mock.Arguments) {
			input := args.Get(2).(entity.UserProfile)
			mockRepo.ExpectedCalls[0].ReturnArguments = mock.Arguments{input, nil}
		}).
		Once()

	creator := userprofile.NewCreator(
		logger,
		mockRepo,
		&userprofile.UserProfileMapper{},
	)

	input := faker.UserProfile()
	input.FirstName = input.FirstName + "  "
	input.LastName = "  " + input.LastName

	inputSanitized := input
	inputSanitized.Sanitize()

	result, err := creator.CreateUserProfileTx(context.Background(), &sql.Tx{}, input)

	mockRepo.AssertNumberOfCalls(t, "InsertUserProfile", 1)
	assert.NilError(t, err)
	assert.DeepEqual(t, inputSanitized, result)
}

func TestCreator_CreateUserProfileTx_ValidationError(t *testing.T) {

	logger, err := logx.InitLogger()
	if err != nil {
		t.Fatalf("failed to initialize logger: %v", err)
	}

	tcc := []struct {
		name    string
		inputFn func(input *userprofile.UserProfile)
	}{
		{
			name: "missing first name",
			inputFn: func(input *userprofile.UserProfile) {
				input.FirstName = "   "
			},
		},
		{
			name: "missing last name",
			inputFn: func(input *userprofile.UserProfile) {
				input.LastName = "   "
			},
		},
		{
			name: "date of birth in the future",
			inputFn: func(input *userprofile.UserProfile) {
				input.DateOfBirth = civil.DateOf(time.Now().Add(time.Hour * 24))
			},
		},
		{
			name: "zero date of birth",
			inputFn: func(input *userprofile.UserProfile) {
				input.DateOfBirth = civil.Date{}
			},
		},
		{
			name: "invalid date of birth",
			inputFn: func(input *userprofile.UserProfile) {
				input.DateOfBirth = civil.Date{Year: 2024, Month: 13, Day: 32}
			},
		},
	}

	for _, tc := range tcc {
		t.Run(tc.name, func(t *testing.T) {
			mockRepo := new(mockCreatorRepo)

			creator := userprofile.NewCreator(
				logger,
				mockRepo,
				&userprofile.UserProfileMapper{},
			)

			input := faker.UserProfile()
			tc.inputFn(&input)

			_, err = creator.CreateUserProfileTx(context.Background(), nil, input)

			mockRepo.AssertNumberOfCalls(t, "InsertUserProfile", 0)
			assert.ErrorType(t, err, &errorx.ValidationError{})
		})
	}
}
