package enum

import (
	"fmt"
)

type RangeDuration int

func (r RangeDuration) String() string {
	switch r {
	case DAY:
		return "DAY"
	case WEEK:
		return "WEEK"
	case MONTH:
		return "MONTH"
	}

	return fmt.Sprintf("%%!RangeDuration(%d)", r)
}

const (
	DAY RangeDuration = iota + 1
	WEEK
	MONTH
)

func NewRangeDurationByString(str string) (RangeDuration, error) {
	switch str {
	case "day":
		return DAY, nil
	case "week":
		return WEEK, nil
	case "month":
		return MONTH, nil
	}

	return 0, fmt.Errorf("unknown range")
}
