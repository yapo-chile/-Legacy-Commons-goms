package repository

import (
	"errors"
	"fmt"
	"math"
	"strconv"
	"strings"

	"github.schibsted.io/Yapo/goms/pkg/domain"
)

// WayStairRepo instantiate all functions needed so can may be able to
// calculate possible ways and its combinations
type WayStairRepo struct{}

// NewWayStair initialice WayStairRepo
func NewWayStair() *WayStairRepo {
	return &WayStairRepo{}
}

// removeDuplicates utility function that takes an slice, loop itself
// and returns a new slice with non duplicated items
func (r *WayStairRepo) removeDuplicates(elements []string) []string {
	// Use map to record duplicates as we find them.
	encountered := map[string]bool{}
	result := []string{}
	for v := range elements {
		if encountered[elements[v]] == true {
			// Do not add duplicate.
		} else {
			// Record this element as an encountered element.
			encountered[elements[v]] = true
			// Append to result slice.
			result = append(result, elements[v])
		}
	}
	return result
}

// padRight utility function that fill with pad whats required by lenght
// Ex. padRight("b", "x", 5) = output: bxxxx
func (r *WayStairRepo) padRight(str, pad string, lenght int) string {
	for {
		str += pad
		if len(str) > lenght {
			return str[0:lenght]
		}
	}
}

// convertNatural takes any number from any base to return it as base 10/decimal
func (r *WayStairRepo) convertNatural(num string, base int) (int, error) {
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
func (r *WayStairRepo) getComb(num int, nstair int) (string, error) {
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
func (r *WayStairRepo) Calculate(nth int) (domain.WayStair, error) {
	var combs []string
	// validates that input is not negative
	if nth < 0 {
		return domain.WayStair{}, fmt.Errorf("out of range")
	}
	// get values of the lowest and max possible ternary number of the
	// requested nth
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
		Ways:  len(response),
		Combs: fmt.Sprintf("{%s}", strings.Join(response, ",")),
	}, nil
}
