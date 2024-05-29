package timePrec

import (
	"github.com/beevik/ntp"
	"time"
)

func GetTimePrec() (time.Time, error) {
	return ntp.Time("0.beevik-ntp.pool.ntp.org")
}
