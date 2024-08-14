package gotesting

import (
	"testing"

	"github.com/stretchr/testify/assert"
	gomock "go.uber.org/mock/gomock"
)

func Test_GetUser_HappyCase(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockedUserServiceClient := NewMockUserServiceClient(ctrl)

	// We are expecting this call to GetUser and should have id 123
	mockedUserServiceClient.EXPECT().GetUser("123").
		Return(&User{ID: "123", Name: "shubham"}, nil)

	authService := NewAuthService(mockedUserServiceClient)
	user, err := authService.GetUser("123")
	assert.NoError(t, err)
	assert.Equal(t, user.Name, "shubham")
}

func Test_GetUser_Failure(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockedUserServiceClient := NewMockUserServiceClient(ctrl)

	// We are expecting this call to GetUser and should have id 1234
	// Which will not match id we will be calling for mock
	// This will fail since no matching mocks were found
	mockedUserServiceClient.EXPECT().GetUser("1234").
		Return(&User{ID: "123", Name: "shubham"}, nil)

	authService := NewAuthService(mockedUserServiceClient)
	user, err := authService.GetUser("123")
	assert.Error(t, err)
	assert.Nil(t, user)
}
