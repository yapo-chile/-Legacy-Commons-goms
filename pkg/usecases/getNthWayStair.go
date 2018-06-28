package usecases

import (
	"fmt"

	"github.schibsted.io/Yapo/goms/pkg/domain"
)

type GetNthWayStairUsecase interface {
	GetNth(n int) (domain.WayStair, error)
}

// WayStairInteractorLogger defines all the events a WayStairInteractor may
// need/like to report as they happen
type WayStairInteractorLogger interface {
	LogBadInput(int)
	LogRepositoryError(int, domain.WayStair, error)
}

// WayStairInteractor implements GetNthWayStairUsecase by using Repository
// to store new WayStair at required and to retrieve the final answer.
type WayStairInteractor struct {
	Logger     WayStairInteractorLogger
	Repository domain.WayStairRepository
}

// GetNth finds the nth WayStair Number by recursively generating one more
// from the last know pair. The step base is O(n).
func (interactor *WayStairInteractor) GetNth(n int) (domain.WayStair, error) {
	// Ensure correct input
	top := 14
	if n <= 0 {
		interactor.Logger.LogBadInput(n)
		return domain.WayStair{}, fmt.Errorf("there's no such thing as %dth WayStair", top)
	}
	if n > top {
		interactor.Logger.LogBadInput(n)
		return domain.WayStair{}, fmt.Errorf("stop eating donnuts, you cant go higher than %d", top)
	}
	// Check if the repository already knows it

	ws, err := interactor.Repository.WayStairs(n)
	if err != nil {
		// Report the error
		interactor.Logger.LogRepositoryError(n, ws, err)
		return domain.WayStair{}, fmt.Errorf("error, this value is impossible to calculate")
	}
	return ws, nil
}
