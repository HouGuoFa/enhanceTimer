/*
	author houguofa
	Copyright 2015 houguofa. All rights reserved.
*/
package timer

import (
	"sync"
	"time"
)

const (
	maxListNode = 20000
	dxyAvg      = interval / 10
)

type list struct {
	header    *node
	nodeCount int
	llock     *sync.RWMutex
}

func createList() *list {
	return &list{header: getHeadNode(), nodeCount: 0, llock: new(sync.RWMutex)}
}

func (l *list) add(t, solt int) *Timer {
	l.llock.Lock()
	defer l.llock.Unlock()
	if l.nodeCount+1 > maxListNode {
		return nil
	}

	l.nodeCount++
	n := getNode(t)
	l.insert(n)

	return n.t
}

func (l *list) append(n *node, solt int) {
	l.llock.Lock()
	defer l.llock.Unlock()

	l.nodeCount++
	l.insert(n)
}

func (l *list) insert(n *node) {
	nptr := l.header

	if nptr.next != nil {
		nptr.next.prec = n
	}

	n.next = nptr.next
	nptr.next = n
	n.prec = nptr
}

func (l *list) onTick(ct int) {

	now := time.Now().Unix() + dxyAvg
	l.llock.RLock()
	nptr := l.header.next
	dnodes := make([]*node, 0, 100)
	cnodes := make([]*node, 0, 100)

	for nptr != nil {

		if !nptr.t.alive {
			dnodes = append(dnodes, nptr)
			nptr = nptr.next
			continue
		}

		if nptr.deadline <= now {
			select {
			case nptr.t.C <- true:
			default:
			}
			cnodes = append(cnodes, nptr)
			nptr = nptr.next
			continue
		}
		nptr = nptr.next
	}
	l.llock.RUnlock()
	l.adjust(ct, dnodes, cnodes, now)
}

func (l *list) adjust(ct int, delnodes, cgnodes []*node, now int64) {

	l.llock.Lock()
	defer l.llock.Unlock()
	l.nodeCount -= len(cgnodes)
	if len(cgnodes) > 0 {
		changeNodes(ct, cgnodes, now)
	}

	l.nodeCount -= len(delnodes)
	if len(delnodes) > 0 {
		freeNodes(delnodes)
	}
}

func changeNodes(ct int, cgnodes []*node, now int64) {
	for _, n := range cgnodes {
		delNode(n)
		n.deadline = now + int64(n.inter)
		solt := (ct + n.inter) % interval
		l := getReadList(solt)
		l.append(n, solt)
	}
}

func freeNodes(dnodes []*node) {
	for _, n := range dnodes {
		delNode(n)
		TimerPool.Put(n.t)
		n.t = nil
		nodePool.Put(n)
	}
}

func delNode(n *node) {
	if n.next != nil {
		n.next.prec = n.prec
	}
	n.prec.next = n.next
	n.prec = nil
	n.next = nil
}

// ----------------------------------------------------------------
//func stat(l *list) {

//	alive := 0
//	stop := 0
//	total := 0
//	n := l.header.next
//	for n != nil {
//		total++
//		if n.t.alive {
//			alive++
//		} else {
//			stop++
//		}
//		n = n.next
//	}
//	golog.Info("%%%%%%%%%%%%%%%---->", total, alive, stop)
//}
