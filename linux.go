// +build !darwin
package hybridtime

import (
	"fmt"

	"golang.org/x/sys/unix"
)

func (c *Clock) walltimeWithError() (nowUsec uint64, errUsec uint64, err error) {
	t := new(unix.Timex)
	status, err := unix.Adjtimex(t)
	if err != nil {
		return 0, 0, nil
	}
	if status != 0 {
		return 0, 0, fmt.Errorf("Got status code %d", status)
	}
	nowUsec = t.Time.Sec*nanosPerSecond + t.Time.Usec/c.divisor
	errUsec = t.Maxerr
	return nowUsec, errUsec, nil
}

func (c *Clock) getDivisor() (uint64, error) {
	t := new(unix.Timex)
	status, err := unix.Adjtimex(t)
	if err != nil {
		return 0, 0, nil
	}
	if status != 0 {
		return 0, 0, fmt.Errorf("Got status code %d", status)
	}
	if (timex.Status & STA_NANO) == STA_NANO {
		return nanosecondDivisor
	} else {
		return microsecondDivisor
	}
}
