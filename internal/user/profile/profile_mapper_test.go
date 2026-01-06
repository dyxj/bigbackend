package profile_test

import (
	"testing"

	"github.com/dyxj/bigbackend/internal/user/profile"
	"github.com/dyxj/bigbackend/test/faker"
	"github.com/stretchr/testify/assert"
)

func TestUserProfileMapper_CreateRequestToModel(t *testing.T) {
	mapper := &profile.UserProfileMapper{}

	createReq := faker.UserProfileCreateRequest()

	model := mapper.CreateRequestToModel(createReq)

	assert.Equal(t, createReq.UserID, model.UserID)
	assert.Equal(t, createReq.FirstName, model.FirstName)
	assert.Equal(t, createReq.LastName, model.LastName)
	assert.Equal(t, createReq.DateOfBirth, model.DateOfBirth)
}

func TestUserProfileMapper_EntityToModel(t *testing.T) {
	mapper := &profile.UserProfileMapper{}

	entity := faker.UserProfileEntity()
	entity.Version = 10

	model := mapper.EntityToModel(entity)

	assert.Equal(t, entity.ID, model.ID)
	assert.Equal(t, entity.UserID, model.UserID)
	assert.Equal(t, entity.FirstName, model.FirstName)
	assert.Equal(t, entity.LastName, model.LastName)
	assert.Equal(t, entity.DateOfBirth, model.DateOfBirth)
	assert.Equal(t, entity.CreateTime, model.CreateTime)
	assert.Equal(t, entity.UpdateTime, model.UpdateTime)
	assert.Equal(t, entity.Version, model.Version)
}

func TestUserProfileMapper_ModelToEntity(t *testing.T) {
	mapper := &profile.UserProfileMapper{}

	userProfile := faker.UserProfile()
	userProfile.Version = 5

	entityProfile := mapper.ModelToEntity(userProfile)

	assert.Equal(t, userProfile.ID, entityProfile.ID)
	assert.Equal(t, userProfile.UserID, entityProfile.UserID)
	assert.Equal(t, userProfile.FirstName, entityProfile.FirstName)
	assert.Equal(t, userProfile.LastName, entityProfile.LastName)
	assert.Equal(t, userProfile.DateOfBirth, entityProfile.DateOfBirth)
	assert.Equal(t, userProfile.CreateTime, entityProfile.CreateTime)
	assert.Equal(t, userProfile.UpdateTime, entityProfile.UpdateTime)
	assert.Equal(t, userProfile.Version, entityProfile.Version)
}

func TestUserProfileMapper_ModelToResponse(t *testing.T) {
	mapper := &profile.UserProfileMapper{}

	userProfile := faker.UserProfile()
	userProfile.Version = 3

	response := mapper.ModelToResponse(userProfile)

	assert.Equal(t, userProfile.ID, response.ID)
	assert.Equal(t, userProfile.UserID, response.UserID)
	assert.Equal(t, userProfile.FirstName, response.FirstName)
	assert.Equal(t, userProfile.LastName, response.LastName)
	assert.Equal(t, userProfile.DateOfBirth, response.DateOfBirth)
	assert.Equal(t, userProfile.CreateTime, response.CreateTime)
	assert.Equal(t, userProfile.UpdateTime, response.UpdateTime)
	assert.Equal(t, userProfile.Version, response.Version)
}
