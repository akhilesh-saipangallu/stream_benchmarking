package utils

import (
	"fmt"
	"sync/atomic"
)

type LogOnceEveryN struct {
	counter uint64
	n       uint64
}

func NewLogOnceEveryN(n uint64) *LogOnceEveryN {
	if n == 0 {
		n = 1
	}
	return &LogOnceEveryN{n: n}
}

func (l *LogOnceEveryN) Log(format string, a ...interface{}) {
	currentCount := atomic.AddUint64(&l.counter, 1)

	if currentCount%l.n == 0 {
		fmt.Printf(format+"\n", a...)
	}
}
