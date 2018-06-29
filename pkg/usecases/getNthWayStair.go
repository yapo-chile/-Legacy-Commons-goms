package usecases

import (
	"fmt"

	"github.schibsted.io/Yapo/goms/pkg/domain"
)

// GetNthWayStairUsecase interface that uses WayStairInteractor
// so it can be used by other layers
type GetNthWayStairUsecase interface {
	GetNth(n int) (domain.WayStair, error)
}

// WayStairInteractorLogger defines all the events as WayStairInteractor may
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

// GetNth finds the nth Ways by using Repository to calculate possible ways
// and its combination to retrieve them as answer.
func (interactor *WayStairInteractor) GetNth(n int) (domain.WayStair, error) {
	// Ensure that input is not negative
	if n <= 0 {
		interactor.Logger.LogBadInput(n)
		return domain.WayStair{}, fmt.Errorf("there's no such thing as %dth WayStair", n)
	}
	// Ensure that input is not that high so it could ran away due memory.
	top := 14
	if n > top {
		interactor.Logger.LogBadInput(n)
		return domain.WayStair{}, fmt.Errorf("stop eating donnuts, you cant go higher than %d", top)
	}
	// Check if the repository already knows it
	ws, err := interactor.Repository.Calculate(n)
	if err != nil {
		// Report the error
		interactor.Logger.LogRepositoryError(n, ws, err)
		return domain.WayStair{}, fmt.Errorf("error, this value is impossible to calculate")
	}
	return ws, nil
}
