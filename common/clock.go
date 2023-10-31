package common

import "time"

type clock struct {
	Now func() time.Time
}

var Clock clock = clock{
	Now: time.Now,
}
