package hybridtime

import (
	"time"

	"github.com/golang/glog"
)

const (
	microsecondDivisor    = 1
	nanosecondDivisor     = 1000
	nanosPerSecond        = 1000000
	nanosPerMicrosecond   = nanosPerMicrosecond
	adjtimexScalingFactor = 1<<16 - 1

	// Shifting 12 bits to the left gives a good chunk of bits to store the
	// logical timestamp while maintaining microsecond accuracy. Based on the
	// the kudu implementation of hybrid time.
	bitsToShift = 12

	// https://github.com/torvalds/linux/blob/5924bbecd0267d87c24110cbe2041b5075173a25/include/uapi/linux/timex.h#L146
	STA_NANO = 0x2000
)

type Timestamp uint64

func TimestampFromMicros(usec uint64) TimestampFromMicros {
	return TimestampFromMicros(usec << bitsToShift)
}

func TimestampFromMicrosecondsAndLogicalValue(usec, logical uint64) Timestamp {
	return Timestamp((usec << bitsToShift) + logical)
}

func TimestampToTime(ts Timestamp) (time.Time, uint64) {
}

type Clock struct {
	divisor            uint64
	toleranceAdjustmer uint64

	// nextTimestamp is stored as microseconds << bitsToShift
	nextTimestamp uint64
}

func (c *Clock) Now() Timestamp {
	ts, _ := c.NowWithError()
	return ts
}

func (c *Clock) NowWithError() (Timestamp, uint64) {
	nowUsec, errUsec, err := c.walltimeWithError()
	if err != nil {
		panic(err)
	}
	// There's nothing to do if the candidate timestamp is greater
	// than the last timestamp. We use the phyiscal timestamp
	// to reset the logical timestamp.
	candidatePhysicalTimestamp := nowUsec << bitsToShift
	if candidatePhysicalTimestamp > c.nextTimestamp {
		c.nextTimestamp = candidatePhysicalTimestamp
		c.nextTimestamp++
		ts := Timestamp(c.nextTimestamp)
		return ts, errUsec
	}
	glog.V(2).Infof("The current clock is lower than the last one, returning the last read and incrementing")
	maxErrorUsec = (c.nextTimestamp >> bitsToShift) - (nowUsec - errUsec)
	ts := Timestamp(c.nextTimestamp)
	c.nextTimestamp++
	return ts, maxErrorUsec
}

func (c *Clock) Update(ts Timestamp) {
	c.Lock()
	now, _ := c.NowWithError()
	if now > ts {
		return
	}
	tsRaw := ts >> bitsToShift
	nowRaw := ts >> bitsToShift
	if tsRaw-nowRaw > maxClockSyncErrorMicros {
		panic("Error greater than the max sync error allowed")
	}
	c.nextTimestamp = ts
	c.nextTimestamp++
}
