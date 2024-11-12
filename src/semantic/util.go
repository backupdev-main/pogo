package semantic

import (
	"strconv"
)

func isNumeric(s string) bool {
	_, errInt := strconv.Atoi(s)
	_, errFloat := strconv.ParseFloat(s, 64)
	return errInt == nil || errFloat == nil
}
