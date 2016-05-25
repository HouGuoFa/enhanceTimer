/*
	author houguofa
	Copyright 2015 houguofa. All rights reserved.
*/
package enhanceTimer

import (
	"sync"
	"sync/atomic"
	"time"
)

//var maxUnitTimers int32 = unitSlotCount * unitTaskCount

type unitTimer struct {
	totalTimers  int32
	minders      map[int]*unitNodeMinder
	lock         *sync.RWMutex
	curMaxMinder int
}

func newUnitTimer() *unitTimer {
	ut := &unitTimer{totalTimers: 0,
		minders:      make(map[int]*unitNodeMinder),
		lock:         new(sync.RWMutex),
		curMaxMinder: -1}
	go ut.run()
	return ut
}

func (ut *unitTimer) add() *Timer {

	ut.lock.RLock()
	for idt, nm := range ut.minders {
		if nm == nil {
			ut.lock.RUnlock()
			ut.lock.Lock()
			if nm == nil {
				nm = newMinder()
				ut.minders[idt] = nm
			}
			ut.lock.Unlock()
			ut.lock.RLock()
		}

		t := nm.add(idt)
		if t != nil {
			atomic.AddInt32(&ut.totalTimers, 1)
			ut.lock.RUnlock()
			return t
		}
	}
	ut.lock.RUnlock()

	return ut.increaseMinder()

}

func (ut *unitTimer) increaseMinder() *Timer {
	ut.lock.Lock()
	defer ut.lock.Unlock()

	ut.curMaxMinder++
	nm := newMinder()
	ut.minders[ut.curMaxMinder] = nm
	t := nm.add(ut.curMaxMinder)
	if t != nil {
		atomic.AddInt32(&ut.totalTimers, 1)
	}
	return t
}

func (ut *unitTimer) del(idt index) {
	minder := ut.minders[idt.firIndex]
	minder.del(idt.secIndex)
	atomic.AddInt32(&ut.totalTimers, -1)
}

func (ut *unitTimer) run() {
	tick := time.Tick(1 * time.Second)
	count := 0
	for {
		select {
		case <-tick:
			//			golog.Info("unit timer minder[", len(ut.minders), "]", ut.curMaxMinder, ut.totalTimers)
			for _, m := range ut.minders {
				if m != nil {
					go m.tick()
				}
			}
			count++
			if count >= 60 {
				//golog.Info("unit timer minders[", len(ut.minders), "]", ut.curMaxMinder, ut.totalTimers)
				count = 0
			}
		}
	}
}

// ----------------------------------------------------------
type unitNodeMinder struct {
	totalNodes int
	full       bool
	empty      bool
	minderLock *sync.RWMutex
	solts      map[int]*Timer
	curMaxSolt int
}

func newMinder() *unitNodeMinder {
	nm := &unitNodeMinder{totalNodes: 0, full: false, empty: true,
		minderLock: new(sync.RWMutex), solts: make(map[int]*Timer), curMaxSolt: -1}
	return nm
}

func (nm *unitNodeMinder) add(solt int) *Timer {
	nm.minderLock.Lock()
	addFlag := false
	defer func() {
		if nm.empty && addFlag {
			nm.empty = false
		}

		if addFlag {
			nm.totalNodes++
		}

		if nm.totalNodes >= unitTasksPreSolt {
			nm.full = true
		}
		nm.minderLock.Unlock()
	}()

	if nm.full {
		return nil
	}

	//t := getTimer(unitScheduler)
	t := TimerPool.Get().(*Timer)
	t.site.timerTeam = unitScheduler
	t.site.firIndex = solt

	for idx, tm := range nm.solts {
		if tm == nil {
			t.site.secIndex = idx
			nm.solts[idx] = t
			addFlag = true
			return t
		}
	}

	nm.curMaxSolt++
	if nm.curMaxSolt <= unitTasksPreSolt {
		t.site.secIndex = nm.curMaxSolt
		nm.solts[nm.curMaxSolt] = t
		addFlag = true
		return t
	}
	return nil
}

func (nm *unitNodeMinder) del(solt int) {

	nm.minderLock.Lock()
	defer nm.minderLock.Unlock()
	t := nm.solts[solt]
	if t != nil {
		TimerPool.Put(t)
	}

	nm.solts[solt] = nil
	nm.totalNodes--
	if nm.totalNodes <= 0 {
		nm.empty = true
	}
	if nm.full {
		nm.full = false
	}
}

func (nm *unitNodeMinder) tick() {
	nm.minderLock.RLock()
	defer nm.minderLock.RUnlock()

	for _, t := range nm.solts {
		if t != nil {
			select {
			case t.C <- true:
			default:
			}
		}
	}
}
