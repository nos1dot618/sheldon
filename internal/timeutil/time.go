package timeutil

import (
	"time"

	"github.com/olebedev/when"
	"github.com/olebedev/when/rules/en"
)

func ParseTime(input string) (time.Time, error) {
	w := when.New(nil)
	w.Add(en.All...)

	result, err := w.Parse(input, time.Now())
	if err != nil || result == nil {
		return time.Time{}, err
	}

	return result.Time, nil
}
