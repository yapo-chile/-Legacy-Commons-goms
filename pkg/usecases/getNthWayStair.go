package usecases

import (
	"errors"
	"fmt"
	"math"
	"strconv"
	"strings"

	"github.schibsted.io/Yapo/goms/pkg/domain"
)

// removeDuplicates utility function that takes an slice, loop itself
// and returns a new slice with non duplicated items
func (r *WayStairInteractor) removeDuplicates(elements []string) []string {
	// Use map to record duplicates as we find them
	encountered := map[string]bool{}
	result := []string{}
	for _, v := range elements {
		if !encountered[v] {
			// Record this element as an encountered element.
			encountered[v] = true
			// Append to result slice.
			result = append(result, v)
		}
	}
	return result
}

// padRight utility function that fill with pad whats required by length
// Ex. padRight("b", "x", 5) = output: bxxxx
func (r *WayStairInteractor) padRight(str, pad string, length int) string {
	for i := length; len(str) <= length; i++ {
		str += pad
	}
	return str
}

// convertNatural takes any number from any base to return it as base 10/decimal
func (r *WayStairInteractor) convertNatural(num string, base int) (int, error) {
	response := 0
	for i := range string(num) {
		strNum := string(num[i])
		floatNum, err := strconv.ParseFloat(strNum, 64)
		if err != nil {
			return 0, errors.New("Error")
		}
		exp := float64(len(num) - i - 1)
		b := float64(base)
		response = int(floatNum*math.Pow(b, exp)) + response
	}
	return int(response), nil
}

// getComb takes a base 10 number to transform it
// to base 3, split them into a slice and checks if the sum of its items are equal to nstair
// returning a comb and if its valid or not
func (r *WayStairInteractor) getComb(num int, nstair int) (string, error) {
	num64 := int64(num)
	nstair64 := int64(nstair)
	// takes the number form base10 to base 3 and split it into an slice
	slice := strings.Split(strconv.FormatInt(num64, 3), "")
	var sum int64
	// loop to store the total sum of the newly created slice
	for o := range slice {
		n, _ := strconv.ParseInt(slice[o], 10, 64)
		sum += n
	}
	// transforns slice into a string and removes zeros
	response := strings.Join(slice, "")
	response = strings.Replace(response, "0", "", -1)
	// if the total sum is equal to the nth requested stair
	// returns nil as well
	if sum == nstair64 {
		return response, nil
	}
	// else returns with error
	return response, errors.New("Bad Comb")
}

// Calculate its where magic happens, if this function were Sinatra this
// would be his New York, it takes an nth, and calculate possible ways and
// its combinations, also checks if there is any repeated item to remove
func (r *WayStairInteractor) Calculate(nth int) (domain.WayStair, error) {
	var combs []string
	// get values of the lowest and max possible ternary number of the
	// requested nth
	// Ensure that input is not negative
	if nth <= 0 {
		r.Logger.LogBadInput(nth)
		return domain.WayStair{},
			fmt.Errorf("there's no such thing as %dth WayStair", nth)
	}
	// Ensure that input is not that high so it could ran away due memory.
	if nth > r.StairsLimit {
		r.Logger.LogBadInput(nth)
		return domain.WayStair{},
			fmt.Errorf("stop eating donnuts, you cant go higher than %d", r.StairsLimit)
	}
	startOn, _ := r.convertNatural(r.padRight("1", "0", nth), 3)
	finishOn, _ := r.convertNatural(r.padRight("2", "2", nth), 3)
	for i := startOn; i < finishOn; i++ {
		comb, e := r.getComb(i, nth)
		if e == nil {
			combs = append(combs, comb)
		}
	}
	// removes duplicates from final slice
	response := r.removeDuplicates(combs)
	return domain.WayStair{
		Stair: domain.Stair(nth),
		Ways:  domain.Ways(len(response)),
		Combs: domain.Combs(fmt.Sprintf("{%s}", strings.Join(response, ","))),
	}, nil
}

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
	Logger      WayStairInteractorLogger
	Repository  domain.WayStairRepository
	StairsLimit int
}

// GetNth finds the nth Ways by using Repository to calculate possible ways
// and its combination to retrieve them as answer.
func (interactor *WayStairInteractor) GetNth(n int) (domain.WayStair, error) {
	wsRepo, err := interactor.Repository.Get(n)
	if err == nil {
		return wsRepo, nil
	}
	interactor.Logger.LogRepositoryError(n, wsRepo, err)
	// Check if the repository already knows it
	ws, err := interactor.Calculate(n)
	if err != nil {
		return domain.WayStair{}, err
	}

	if errSave := interactor.Repository.Save(ws); errSave != nil {
		// Report the error
		interactor.Logger.LogRepositoryError(n, ws, err)
		return domain.WayStair{}, fmt.Errorf("error, can not save the value")
	}

	return ws, nil
}
