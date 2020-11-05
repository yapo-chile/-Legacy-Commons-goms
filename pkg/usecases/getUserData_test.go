package usecases

import (
	"context"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockUserProfileRepository struct {
	mock.Mock
}

func (m *MockUserProfileRepository) GetUserProfileData(ctx context.Context, mail string) (UserBasicData, error) {
	args := m.Called(ctx, mail)
	return args.Get(0).(UserBasicData), args.Error(1)
}

func TestGetUserOk(t *testing.T) {
	m := MockUserProfileRepository{}
	var userb UserBasicData
	ctx := mock.AnythingOfType("*context.emptyCtx")
	m.On("GetUserProfileData", ctx, "").Return(userb, nil)

	i := GetUserDataInteractor{
		UserProfileRepository: &m,
	}
	expected := UserBasicData{"", "", "", "", "", ""}
	output, err := i.GetUser(context.Background(), "")
	assert.NoError(t, err)
	assert.Equal(t, expected, output)
	m.AssertExpectations(t)
}
func TestGetUserError(t *testing.T) {
	m := MockUserProfileRepository{}
	var userb UserBasicData
	ctx := mock.AnythingOfType("*context.emptyCtx")
	m.On("GetUserProfileData", ctx, "").Return(userb, fmt.Errorf("error"))

	i := GetUserDataInteractor{
		UserProfileRepository: &m,
	}

	output, err := i.GetUser(context.Background(), "")
	assert.Error(t, err)
	assert.Empty(t, output)
	m.AssertExpectations(t)
}
