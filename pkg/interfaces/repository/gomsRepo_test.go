package repository

import (
	"context"
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
	mResponse := MockResponse{}
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

	mResponse.On("GetBodyString").Return(string(jsonResponse)).Once()

	ctx := mock.AnythingOfType("*context.emptyCtx")
	mHandler.On("NewRequest", ctx).Return(&mRequest, nil).Once()
	mHandler.On("Send", &mRequest).Return(&mResponse, nil).Once()

	result, err := gomsRepo.GetHealthcheck(context.Background())

	assert.Equal(t, "OK", result)
	assert.Nil(t, err)
	mHandler.AssertExpectations(t)
	mRequest.AssertExpectations(t)
}

func TestGetError(t *testing.T) {
	mHandler := MockHTTPHandler{}
	mResponse := MockResponse{}
	mRequest := MockRequest{}

	gomsRepo := GomsRepository{
		Handler: &mHandler,
	}

	mRequest.On("SetMethod", "GET").Return(&mRequest).Once()
	mRequest.On("SetPath", "").Return(&mRequest).Once()
	mRequest.On("SetTimeOut", mock.AnythingOfType("int")).Return(&mRequest).Once()

	ctx := mock.AnythingOfType("*context.emptyCtx")
	mHandler.On("NewRequest", ctx).Return(&mRequest, nil).Once()
	mHandler.On("Send", &mRequest).Return(&mResponse, errors.New("Error")).Once()

	result, err := gomsRepo.GetHealthcheck(context.Background())

	assert.Equal(t, "", result)
	assert.Error(t, err)
	mHandler.AssertExpectations(t)
	mRequest.AssertExpectations(t)
}

func TestPostParseError(t *testing.T) {
	mHandler := MockHTTPHandler{}
	mResponse := MockResponse{}
	mRequest := MockRequest{}

	gomsRepo := GomsRepository{
		Handler: &mHandler,
	}

	mRequest.On("SetMethod", "GET").Return(&mRequest).Once()
	mRequest.On("SetPath", "").Return(&mRequest).Once()
	mRequest.On("SetTimeOut", mock.AnythingOfType("int")).Return(&mRequest).Once()

	mResponse.On("GetBodyString").Return("").Once()

	ctx := mock.AnythingOfType("*context.emptyCtx")
	mHandler.On("NewRequest", ctx).Return(&mRequest, nil).Once()
	mHandler.On("Send", &mRequest).Return(&mResponse, nil).Once()

	result, err := gomsRepo.GetHealthcheck(context.Background())

	assert.Equal(t, result, "")
	assert.Error(t, err)
	mHandler.AssertExpectations(t)
	mRequest.AssertExpectations(t)
}

func TestGetEmptyResponseError(t *testing.T) {
	mHandler := MockHTTPHandler{}
	mResponse := MockResponse{}
	mRequest := MockRequest{}

	gomsRepo := GomsRepository{
		Handler: &mHandler,
	}

	mRequest.On("SetMethod", "GET").Return(&mRequest).Once()
	mRequest.On("SetPath", "").Return(&mRequest).Once()
	mRequest.On("SetTimeOut", mock.AnythingOfType("int")).Return(&mRequest).Once()

	mResponse.On("GetBodyString").Return("").Once()

	ctx := mock.AnythingOfType("*context.emptyCtx")
	mHandler.On("NewRequest", ctx).Return(&mRequest, nil).Once()
	mHandler.On("Send", &mRequest).Return(&mResponse, nil).Once()

	result, err := gomsRepo.GetHealthcheck(context.Background())

	assert.Equal(t, result, "")
	assert.Error(t, err)
	mHandler.AssertExpectations(t)
	mRequest.AssertExpectations(t)
}
