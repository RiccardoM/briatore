package types

import (
	"fmt"
	"strings"
)

type Output byte

func (o Output) String() string {
	switch o {
	case OutText:
		return "text"

	case OutJSON:
		return "json"

	case OutCSV:
		return "csv"

	default:
		panic(fmt.Errorf("invalid output type: %d", o))
	}
}

const (
	OutText Output = 1
	OutJSON Output = 2
	OutCSV  Output = 3
)

func ParseOutput(out string) (Output, error) {
	switch strings.ToLower(out) {
	case "csv":
		return OutCSV, nil
	case "json":
		return OutJSON, nil
	case "text":
		return OutText, nil
	default:
		return 0, fmt.Errorf("invalid output type: %s", out)
	}
}
