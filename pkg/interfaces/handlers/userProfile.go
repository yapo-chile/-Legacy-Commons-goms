package handlers

import (
	"net/http"

	"github.com/Yapo/goutils"
	"github.mpi-internal.com/Yapo/goms/pkg/usecases"
)

// UserProfileHandler implements the handler interface and responds to
// userProfile requests using an interactor.
type UserProfileHandler struct {
	Interactor usecases.UserProfileInteractor
}

type userProfileRequestInput struct {
	Mail string `json:"Mail"`
}

//userProfileRequestOutput specifies the format of the handler output
type userProfileRequestOutput struct {
	Name    string `json:"Full name"`
	Phone   string `json:"Cellphone"`
	Gender  string `json:"Gender"`
	Country string `json:"Country"`
	Region  string `json:"Region"`
	Commune string `json:"Commune"`
}

type userProfileRequestError goutils.GenericError

// Input returns a fresh, empty instance of userProfileRequestInput
func (h *UserProfileHandler) Input(ir InputRequest) HandlerInput {
	input := userProfileRequestInput{}
	ir.Set(&input).
		FromJSONBody()
	return &input
}

//Execute is the main function of the userProfile handler
func (h *UserProfileHandler) Execute(ig InputGetter) *goutils.Response {
	input, response := ig()
	if response != nil {
		return response
	}
	in := input.(*userProfileRequestInput)

	userBasicData, err := h.Interactor.GetUser(in.Mail)

	if err != nil {
		return &goutils.Response{
			Code: http.StatusNoContent,
		}
	}

	return &goutils.Response{
		Body: h.fillInternalOutput(userBasicData),
		Code: http.StatusOK,
	}
}

// fillInternalOutput parses the data retrieved from the handler to handler expected output
func (h *UserProfileHandler) fillInternalOutput(userBasicData usecases.UserBasicData) userProfileRequestOutput {

	return userProfileRequestOutput{
		Name:    userBasicData.Name,
		Phone:   userBasicData.Phone,
		Gender:  userBasicData.Gender,
		Country: userBasicData.Country,
		Region:  userBasicData.Region,
		Commune: userBasicData.Commune,
	}
}
