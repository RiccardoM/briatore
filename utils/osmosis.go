package utils

import (
	"fmt"
	"regexp"
	"strconv"
)

var (
	poolCoinDenomRegex = regexp.MustCompile(`\d+`)
)

func ParsePoolID(tokenDenom string) (uint64, error) {
	matches := poolCoinDenomRegex.FindStringSubmatch(tokenDenom)
	if len(matches) == 0 {
		return 0, fmt.Errorf("no number found in string")
	}
	return strconv.ParseUint(matches[0], 10, 64)
}
