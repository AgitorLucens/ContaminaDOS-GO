package utils

import (
	"fmt"
	"strconv"
)

func ToInt(s string) (int, error) {
	i, err := strconv.Atoi(s)
	if err != nil {
		return 0, err
	}
	return i, nil
}


func ToString(i any) string {
	return fmt.Sprintf("%s", i)
}