package handlers

import (
	"net/http"

	"github.com/Yapo/goutils"
	"github.schibsted.io/Yapo/goms/pkg/usecases"
)

// WayStairHandler implements the handler interface and responds to
// /WayStair requests using an interactor. It's purpose is just to
// demonstrate Clean Architecture with a practical scenario

type WayStairHandler struct {
	Interactor usecases.GetNthWayStairUsecase
}
type wayStairRequestInput struct {
	N int `json:"n"`
}

type wayStairRequestOutput struct {
	Ways         int    `json:"ways"`
	Combinations string `json:"combinations"`
}

type wayStairRequestError goutils.GenericError

// Input returns a fresh, empty instance of wayStairRequestInput
func (h *WayStairHandler) Input() HandlerInput {
	return &wayStairRequestInput{}
}

// Execute carries on a /waystair request. Uses the given interactor to carry out
// the operation and get the desired value. Expected body format:
//	{
// 		Ways: int - number of possible combinations
//		Combs: String - all possible combinations
//	}
// Expected response format:
//   { Ways: int - Operation result,
//     Combs: slice combs }
// Expected error format:
//   { ErrorMessage: string - Error detail }
func (h *WayStairHandler) Execute(ig InputGetter) *goutils.Response {
	input, response := ig()
	if response != nil {
		return response
	}

	in := input.(*wayStairRequestInput)
	f, err := h.Interactor.GetNth(in.N)
	if err != nil {
		return &goutils.Response{
			Code: http.StatusBadRequest,
			Body: wayStairRequestError{
				err.Error(),
			},
		}
	}
	return &goutils.Response{
		Code: http.StatusOK,
		Body: wayStairRequestOutput{
			Ways:         int(f.Ways),
			Combinations: string(f.Combs),
		},
	}
}
