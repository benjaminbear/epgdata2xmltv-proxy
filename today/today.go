package today

import (
	"strconv"
	"time"
)

const DayFormat = "20060102"

type Today struct {
	time.Time
}

func New() *Today {
	return &Today{time.Now()}
}

func (t *Today) GetDayPlus(days int) string {
	return t.AddDate(0, 0, days).Format(DayFormat)
}

func (t *Today) GetString() string {
	return t.Format(DayFormat)
}

func (t *Today) GetInt() (int, error) {
	return strconv.Atoi(t.Format(DayFormat))
}
