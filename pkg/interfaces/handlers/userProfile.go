package handlers

import (
	"net/http"

	"github.com/Yapo/goutils"
	"github.mpi-internal.com/Yapo/goms/pkg/usecases"
)

// UserProfileHandler implements the handler interface and responds to
// /userProfile requests using an interactor. It's purpose is just to
// demonstrate Clean Architecture with a practical scenario
type UserProfileHandler struct {
	Interactor usecases.UserProfileInteractor
}

type userProfileRequestInput struct {
	Mail string `json:"Mail"`
}

type userProfileRequestOutput struct {
	UseraBasicData usecases.UserBasicData
}

type userProfileRequestError goutils.GenericError

// Input returns a fresh, empty instance of userProfileRequestInput
func (h *UserProfileHandler) Input(ir InputRequest) HandlerInput {
	input := userProfileRequestInput{}
	ir.Set(&input).
		FromJSONBody()
	return &input
}

func (h *UserProfileHandler) Execute(ig InputGetter) *goutils.Response {
	input, response := ig()
	if response != nil {
		return response
	}
	in := input.(*userProfileRequestInput)

	userBasic, err := h.Interactor.GetUser(in.Mail)

	if err != nil {
		return &goutils.Response{
			Code: http.StatusNoContent,
		}
	}

	return &goutils.Response{
		Body: userBasic,
		Code: http.StatusOK,
	}
}
