package utils

import (
	"fmt"
	"strconv"
)

func Atoi(s string) int {
	n, err := strconv.Atoi(s)
	if err != nil {
		panic(fmt.Sprintf("invalid int: %s, err: %v", s, err))
	}
	return n
}
