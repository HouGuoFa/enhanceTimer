/*
	author houguofa
	Copyright 2015 houguofa. All rights reserved.
*/
package enhanceTimer

import (
	"sync/atomic"
	//"time"
)

var addTimers int64
var delTimers int64
var addWheelTimers int64
var delWheelTimers int64

type scheduler struct {
	unitMinder  *unitTimer
	timeWheeler *wheeler
}

var scd *scheduler = nil

func init() {
	scd = &scheduler{unitMinder: newUnitTimer(), timeWheeler: newWheeler()}
	//go func() {
	//	tick := time.Tick(10 * time.Second)
	//	for {
	//		select {
	//		case <-tick:
	//			golog.Info("Timer @@@@ ->",
	//				atomic.LoadInt64(&addTimers), atomic.LoadInt64(&delTimers),
	//				atomic.LoadInt64(&addWheelTimers), atomic.LoadInt64(&delWheelTimers))
	//			atomic.StoreInt64(&addTimers, 0)
	//			atomic.StoreInt64(&delTimers, 0)
	//			atomic.StoreInt64(&addWheelTimers, 0)
	//			atomic.StoreInt64(&delWheelTimers, 0)
	//		}
	//	}
	//}()
}

func NewTimer(d int) *Timer {
	if scd == nil {
		return nil
	}
	t := scd.add(d)
	return t
}

func stop(t *Timer) {
	if scd == nil {
		return
	}

	scd.del(t)
}

func (s *scheduler) add(d int) *Timer {
	if d < 1 {
		d = 1
	}
	if d == unitTick {
		atomic.AddInt64(&addTimers, 1)
		return s.unitMinder.add()
	}

	atomic.AddInt64(&addWheelTimers, 1)
	return s.timeWheeler.add(d)
}

func (s *scheduler) del(t *Timer) {
	if t.site.timerTeam == unitScheduler {
		atomic.AddInt64(&delTimers, 1)
		s.unitMinder.del(t.site)
		return
	}

	atomic.AddInt64(&delWheelTimers, 1)
	s.timeWheeler.del(t)
}
