package tonight

import (
	"fmt"
	"time"

	"github.com/karrick/tparse"
)

const (
	dateFormat = "2006-01-02"
)

func parseDate(date string) (time.Time, error) {
	d, err := time.Parse(dateFormat, date)
	if err == nil {
		return d, nil
	}

	d, err = tparse.ParseNow(dateFormat, fmt.Sprintf("now+%s", date))
	if err != nil {
		return time.Time{}, err
	}

	return d, nil
}
