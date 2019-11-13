package handlers

import (
	"net/http"

	"github.com/Yapo/goutils"
	"github.mpi-internal.com/Yapo/goms/pkg/usecases"
)

// UserProfilePrometheusDefaultLogger is the logger used by the handler
type UserProfilePrometheusDefaultLogger interface {
	LogBadRequest(input interface{})
	LogErrorGettingInternalData(err error)
}

// GetUserDataHandler implements the handler interface and responds to
type GetUserDataHandler struct {
	Interactor UserProfileInteractor
	Logger     UserProfilePrometheusDefaultLogger
}

type userProfileRequestInput struct {
	Mail string `query:"mail"`
}

// userProfileRequestOutput specifies the format of the handler output
type userProfileRequestOutput struct {
	Name    string `json:"fullname"`
	Phone   string `json:"cellphone"`
	Gender  string `json:"gender"`
	Country string `json:"country"`
	Region  string `json:"region"`
	Commune string `json:"commune"`
}

// UserProfileInteractor is the interactor used by the handler
type UserProfileInteractor interface {
	GetUser(mail string) (usecases.UserBasicData, error)
}

// Input returns a fresh, empty instance of userProfileRequestInput
func (h *GetUserDataHandler) Input(ir InputRequest) HandlerInput {
	input := userProfileRequestInput{}
	ir.Set(&input).
		FromQuery()
	return &input
}

// Execute is the main function of the userProfile handler
func (h *GetUserDataHandler) Execute(ig InputGetter) *goutils.Response {
	input, response := ig()
	if response != nil {
		h.Logger.LogBadRequest(response)
		return response
	}
	in := input.(*userProfileRequestInput)

	userBasicData, err := h.Interactor.GetUser(in.Mail)

	if err != nil {
		h.Logger.LogErrorGettingInternalData(err)
		return &goutils.Response{
			Code: http.StatusNoContent,
		}
	}

	return &goutils.Response{
		Code: http.StatusOK,
		Body: h.fillInternalOutput(userBasicData),
	}
}

// fillInternalOutput parses the data retrieved from the handler to the expected output
func (h *GetUserDataHandler) fillInternalOutput(userBasicData usecases.UserBasicData) userProfileRequestOutput {
	return userProfileRequestOutput{
		Name:    userBasicData.Name,
		Phone:   userBasicData.Phone,
		Gender:  userBasicData.Gender,
		Country: userBasicData.Country,
		Region:  userBasicData.Region,
		Commune: userBasicData.Commune,
	}
}
