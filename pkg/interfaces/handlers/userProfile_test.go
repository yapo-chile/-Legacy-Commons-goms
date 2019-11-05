package handlers

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/Yapo/goutils"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.mpi-internal.com/Yapo/goms/pkg/usecases"
)

type mockUserProfilePrometheusDefaultLogger struct {
	mock.Mock
}

func (m *mockUserProfilePrometheusDefaultLogger) LogBadRequest(input interface{}) {
	m.Called(input)
}

func (m *mockUserProfilePrometheusDefaultLogger) LogErrorGettingInternalData(err error) {
	m.Called(err)
}

type mockUserProfileInteractor struct {
	mock.Mock
}

func (m *mockUserProfileInteractor) GetUser(mail string) (usecases.UserBasicData, error) {
	args := m.Called(mail)
	return args.Get(0).(usecases.UserBasicData), args.Error(1)
}
func TestUserProfileHandlerInput(t *testing.T) {
	m := mockUserProfileInteractor{}
	mMockInputRequest := MockInputRequest{}
	mMockInputRequest.On("Set", mock.AnythingOfType("*handlers.userProfileRequestInput")).Return(&mMockInputRequest)
	mMockInputRequest.On("FromJSONBody").Return(&mMockInputRequest)

	h := UserProfileHandler{
		Interactor: &m,
	}
	input := h.Input(&mMockInputRequest)

	var expected *userProfileRequestInput
	assert.IsType(t, expected, input)
	m.AssertExpectations(t)
}
func TestUserProfileHandlerDataRunOK(t *testing.T) {
	mInteractor := &mockUserProfileInteractor{}
	var userb usecases.UserBasicData
	mInteractor.On("GetUser", "").Return(userb, nil)
	h := UserProfileHandler{
		Interactor: mInteractor,
	}
	input := &userProfileRequestInput{
		Mail: "",
	}
	getter := MakeMockInputGetter(input, nil)
	r := h.Execute(getter)

	expected := &goutils.Response{
		Code: http.StatusOK,
		Body: userProfileRequestOutput{
			Name:    "",
			Phone:   "",
			Gender:  "",
			Country: "",
			Region:  "",
			Commune: "",
		},
	}
	assert.Equal(t, expected, r)
	mInteractor.AssertExpectations(t)
}

func TestInternalUserDataHandlerForInternalDataRunError(t *testing.T) {
	mInteractor := &mockUserProfileInteractor{}
	mLogger := &mockUserProfilePrometheusDefaultLogger{}
	err := fmt.Errorf("err")
	var userb usecases.UserBasicData

	mInteractor.On("GetUser", "").Return(userb, err)
	mLogger.On("LogErrorGettingInternalData", err).Once()

	h := UserProfileHandler{
		Interactor: mInteractor,
		Logger:     mLogger,
	}
	input := &userProfileRequestInput{
		Mail: "",
	}
	getter := MakeMockInputGetter(input, nil)
	r := h.Execute(getter)

	expected := &goutils.Response{
		Code: http.StatusNoContent,
	}
	assert.Equal(t, expected, r)
	mLogger.AssertExpectations(t)
	mInteractor.AssertExpectations(t)
}

func TestInternalUserDataHandlerForInternalDataBadRequest(t *testing.T) {
	mInteractor := &mockUserProfileInteractor{}
	mLogger := &mockUserProfilePrometheusDefaultLogger{}

	mLogger.On("LogBadRequest", mock.AnythingOfType("*goutils.Response")).Once()

	h := UserProfileHandler{
		Interactor: mInteractor,
		Logger:     mLogger,
	}
	input := &userProfileRequestInput{
		Mail: "",
	}
	getter := MakeMockInputGetter(input, &goutils.Response{
		Code: http.StatusBadRequest,
	})
	r := h.Execute(getter)

	expected := &goutils.Response{
		Code: http.StatusBadRequest,
	}
	assert.Equal(t, expected, r)
	mLogger.AssertExpectations(t)
	mInteractor.AssertExpectations(t)
}
