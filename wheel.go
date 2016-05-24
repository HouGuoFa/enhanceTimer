/*
	author houguofa
	Copyright 2015 houguofa. All rights reserved.
*/
package timer

import (
	//"qqlive/golog"
	"sync"
	"time"
)

const (
	interval = 256
)

var nodePool = &sync.Pool{New: func() interface{} { return newNode() }}

type node struct {
	t        *Timer
	prec     *node
	next     *node
	inter    int
	deadline int64
}

func newNode() *node {
	return &node{t: nil, prec: nil, next: nil, inter: 2, deadline: 0}
}

type wheeler struct {
	point int
	lock  *sync.RWMutex
	solts []*list
}

var dwheeler *wheeler = nil

func newWheeler() *wheeler {
	w := &wheeler{point: 0, lock: new(sync.RWMutex), solts: make([]*list, interval, interval)}

	for i := 0; i < len(w.solts); i++ {
		w.solts[i] = createList()
	}

	dwheeler = w
	go w.run()
	return w
}

func (w *wheeler) run() {
	tick := time.Tick(1 * time.Second)
	for {
		select {
		case <-tick:
			w.onTick()
		}
	}
}

func (w *wheeler) onTick() {
	w.lock.Lock()
	w.point++
	if w.point >= interval {
		w.point = 0
	}
	w.lock.Unlock()

	go w.solts[w.point].onTick(w.point)
}

func recycleNode(nd *node) {
	delNode(nd)
	TimerPool.Put(nd.t)
	nd.next = nil
	nd.t = nil
	nd.prec = nil
	nodePool.Put(nd)
}

//-------------------------------------------------------

func (w *wheeler) add(t int) *Timer {

	w.lock.RLock()
	solt := (w.point + t + 1) % interval
	w.lock.RUnlock()
	return w.solts[solt].add(t, solt)

}

func getReadList(id int) *list {
	return dwheeler.solts[id]
}

func getNode(t int) *node {
	nd := nodePool.Get().(*node)
	nd.inter = t
	nd.deadline = getDeadLine(t)
	nd.t = TimerPool.Get().(*Timer)
	nd.t.site.timerTeam = wheelScheduler
	nd.t.site.firIndex = secLevel
	nd.t.alive = true
	return nd
}

func getDeadLine(t int) int64 {
	return time.Now().Add(time.Second * time.Duration(int64(t))).Unix()

}
func getHeadNode() *node {
	return &node{nil, nil, nil, -1, -1}
}

func (w *wheeler) del(t *Timer) {
	t.alive = false
}
