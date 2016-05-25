/*
	author houguofa
	Copyright 2015 houguofa. All rights reserved.
*/
package enhanceTimer

import (
	"sync"
	"time"
)

var TimerPool = &sync.Pool{New: func() interface{} { return getTimer(unitScheduler) }}

type index struct {
	timerTeam int
	firIndex  int
	secIndex  int
	thrIndex  int
}

type Timer struct {
	site  index
	C     chan bool
	alive bool
	bkt   *time.Timer
}

func getTimer(nodeType int) *Timer {
	return &Timer{site: index{timerTeam: nodeType},
		C: make(chan bool), alive: true, bkt: nil}
}

func (t *Timer) Stop() {
	stop(t)
}
