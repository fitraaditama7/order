package order

import "time"

type Param struct {
	Limit     int
	Page      int
	Keyword   string
	DateStart *time.Time
	DateEnd   *time.Time
}
