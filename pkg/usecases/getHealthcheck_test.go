package usecases

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockGomsRepository struct {
	mock.Mock
}

func (m *MockGomsRepository) GetHealthcheck(ctx context.Context) (string, error) {
	ret := m.Called(ctx)
	return ret.String(0), ret.Error(1)
}

type MockGomsLogger struct {
	mock.Mock
}

func (m *MockGomsLogger) LogURI(s string) {
	m.Called(s)
}

func (m *MockGomsLogger) LogRequestErr(e error) {
	m.Called(e)
}

func (m *MockGomsLogger) LogHealthcheckOK(s string) {
	m.Called(s)
}

func TestGetHealthcheckOK(t *testing.T) {
	mLogger := MockGomsLogger{}
	mRepo := MockGomsRepository{}

	mInteractor := GetHealthcheckInteractor{
		Logger:         &mLogger,
		GomsRepository: &mRepo,
	}

	mLogger.On("LogURI", mock.AnythingOfType("string"))
	mLogger.On("LogHealthcheckOK", mock.AnythingOfType("string"))

	ctx := mock.AnythingOfType("*context.emptyCtx")
	mRepo.On("GetHealthcheck", ctx).Return("OK", nil)

	resp, err := mInteractor.GetHealthcheck(context.Background())

	assert.Nil(t, err)
	assert.Equal(t, "OK", resp)
	mLogger.AssertExpectations(t)
	mRepo.AssertExpectations(t)
}

func TestGetHealthcheckError(t *testing.T) {
	mLogger := MockGomsLogger{}
	mRepo := MockGomsRepository{}

	mInteractor := GetHealthcheckInteractor{
		Logger:         &mLogger,
		GomsRepository: &mRepo,
	}
	err := errors.New("error")

	mLogger.On("LogURI", mock.AnythingOfType("string"))
	mLogger.On("LogRequestErr", err)

	ctx := mock.AnythingOfType("*context.emptyCtx")
	mRepo.On("GetHealthcheck", ctx).Return("", err)

	resp, err := mInteractor.GetHealthcheck(context.Background())

	assert.Error(t, err)
	assert.Equal(t, "", resp)
	mLogger.AssertExpectations(t)
	mRepo.AssertExpectations(t)
}
