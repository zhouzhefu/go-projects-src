package gotest

import (
	// "fmt"
	"errors"
)

/* 
* Target func to be tested. 
*/
func Division(a, b float64) (float64, error) {
	if b == 0 {
		return 0, errors.New("Cannot divided by 0!")
	}

	return a / b, nil
}