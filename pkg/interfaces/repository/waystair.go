package repository

import (
	"errors"
	"fmt"
	"math"
	"strconv"
	"strings"

	"github.schibsted.io/Yapo/goms/pkg/domain"
)

type WayStairRepo struct {
}

func NewWayStair() *WayStairRepo {
	return &WayStairRepo{}
}

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

func (r *WayStairRepo) padRight(str, pad string, lenght int) string {
	for {
		str += pad
		if len(str) > lenght {
			return str[0:lenght]
		}
	}
}

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

func (r *WayStairRepo) getComb(num int, nstair int) (string, error) {
	num64 := int64(num)
	nstair64 := int64(nstair)
	slice := strings.Split(strconv.FormatInt(num64, 3), "")
	var sum int64
	for o := range slice {
		n, _ := strconv.ParseInt(slice[o], 10, 64)
		sum += n
	}
	response := strings.Join(slice, "")
	response = strings.Replace(response, "0", "", -1)
	if sum == nstair64 {
		return response, nil
	}
	return response, errors.New("Bad Comb")
}

func (r *WayStairRepo) WayStairs(nth int) (domain.WayStair, error) {
	var combs []string
	if nth < 0 {
		return domain.WayStair{}, fmt.Errorf("out of range")
	}
	startOn, _ := r.convertNatural(r.padRight("1", "0", nth), 3)
	finishOn, _ := r.convertNatural(r.padRight("2", "2", nth), 3)
	for i := startOn; i < finishOn; i++ {
		comb, e := r.getComb(i, nth)
		if e == nil {
			combs = append(combs, comb)
		}
	}
	response := r.removeDuplicates(combs)
	return domain.WayStair{
		Ways:  len(response),
		Combs: fmt.Sprintf("{%s}", strings.Join(response, ",")),
	}, nil
}
