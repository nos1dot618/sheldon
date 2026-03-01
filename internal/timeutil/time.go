package timeutil

import (
	"time"

	"github.com/olebedev/when"
	"github.com/olebedev/when/rules/en"
)

func ParseTime(input string) (time.Time, error) {
	w := when.New(nil)
	w.Add(en.All...)

	result, err := w.Parse("2 days ago", time.Now())
	if err != nil || result == nil {
		return time.Now(), err
	}

	return result.Time, nil
}
