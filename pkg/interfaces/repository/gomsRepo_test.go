package repository

import (
	"encoding/json"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestNewGomsRepository(t *testing.T) {
	m := MockHTTPHandler{}

	s := ""
	expected := GomsRepository{
		Handler: &m,
		Path:    s,
		TimeOut: 40,
	}

	result := NewGomsRepository(&m, 40, s)

	assert.Equal(t, &expected, result)
	m.AssertExpectations(t)
}

func TestGetOK(t *testing.T) {
	mHandler := MockHTTPHandler{}
	mRequest := MockRequest{}

	gomsRepo := GomsRepository{
		Handler: &mHandler,
	}
	response := GomsResponse{
		Status: "OK",
	}
	jsonResponse, _ := json.Marshal(response)

	mRequest.On("SetMethod", "GET").Return(&mRequest).Once()
	mRequest.On("SetPath", "").Return(&mRequest).Once()
	mRequest.On("SetTimeOut", mock.AnythingOfType("int")).Return(&mRequest).Once()

	mHandler.On("NewRequest").Return(&mRequest, nil).Once()
	mHandler.On("Send", &mRequest).Return(string(jsonResponse), nil).Once()

	result, err := gomsRepo.GetHealthcheck()

	assert.Equal(t, "OK", result)
	assert.Nil(t, err)
	mHandler.AssertExpectations(t)
	mRequest.AssertExpectations(t)
}

func TestGetError(t *testing.T) {
	mHandler := MockHTTPHandler{}
	mRequest := MockRequest{}

	gomsRepo := GomsRepository{
		Handler: &mHandler,
	}

	mRequest.On("SetMethod", "GET").Return(&mRequest).Once()
	mRequest.On("SetPath", "").Return(&mRequest).Once()
	mRequest.On("SetTimeOut", mock.AnythingOfType("int")).Return(&mRequest).Once()

	mHandler.On("NewRequest").Return(&mRequest, nil).Once()
	mHandler.On("Send", &mRequest).Return(nil, errors.New("Error")).Once()

	result, err := gomsRepo.GetHealthcheck()

	assert.Equal(t, "", result)
	assert.Error(t, err)
	mHandler.AssertExpectations(t)
	mRequest.AssertExpectations(t)
}

func TestPostParseError(t *testing.T) {
	mHandler := MockHTTPHandler{}
	mRequest := MockRequest{}

	gomsRepo := GomsRepository{
		Handler: &mHandler,
	}

	mRequest.On("SetMethod", "GET").Return(&mRequest).Once()
	mRequest.On("SetPath", "").Return(&mRequest).Once()
	mRequest.On("SetTimeOut", mock.AnythingOfType("int")).Return(&mRequest).Once()

	mHandler.On("NewRequest").Return(&mRequest, nil).Once()
	mHandler.On("Send", &mRequest).Return("", nil).Once()

	result, err := gomsRepo.GetHealthcheck()

	assert.Equal(t, result, "")
	assert.Error(t, err)
	mHandler.AssertExpectations(t)
	mRequest.AssertExpectations(t)
}

func TestGetEmptyResponseError(t *testing.T) {
	mHandler := MockHTTPHandler{}
	mRequest := MockRequest{}

	gomsRepo := GomsRepository{
		Handler: &mHandler,
	}

	mRequest.On("SetMethod", "GET").Return(&mRequest).Once()
	mRequest.On("SetPath", "").Return(&mRequest).Once()
	mRequest.On("SetTimeOut", mock.AnythingOfType("int")).Return(&mRequest).Once()

	mHandler.On("NewRequest").Return(&mRequest, nil).Once()
	mHandler.On("Send", &mRequest).Return("", nil).Once()

	result, err := gomsRepo.GetHealthcheck()

	assert.Equal(t, result, "")
	assert.Error(t, err)
	mHandler.AssertExpectations(t)
	mRequest.AssertExpectations(t)
}
