package hybridtime

import "time"

func (c *Clock) walltimeWithError() (nowUsec uint64, errUsec uint64, err error) {
	return time.Now().UnixNano() / c.divisor, 0, nil
}

func (c *Clock) getDivisor() (uint64, error) {
	return nanosecondDivisor
}
