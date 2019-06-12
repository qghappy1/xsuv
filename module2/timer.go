package module2

import (
	"time"

	"github.com/qghappy1/xsuv/util/log"
)

// one dispatcher per goroutine (goroutine not safe)
type dispatcher struct {
	ChanTimer chan *timer
}

func newDispatcher(l int) *dispatcher {
	disp := new(dispatcher)
	disp.ChanTimer = make(chan *timer, l)
	return disp
}

// Timer
type timer struct {
	t  *time.Timer
	cb func()
}

func (t *timer) Stop() {
	t.t.Stop()
	t.cb = nil
}

func (t *timer) Cb() {
	defer log.ErrorPanic()
	if t.cb != nil {
		t.cb()
		t.cb = nil
	}
}

func (disp *dispatcher) AfterFunc(d time.Duration, cb func()) *timer {
	t := new(timer)
	t.cb = cb
	t.t = time.AfterFunc(d, func() {
		disp.ChanTimer <- t
	})
	return t
}

// Cron
type Cron struct {
	t *timer
}

func (c *Cron) Stop() {
	if c.t != nil {
		c.t.Stop()
	}
}

func (disp *dispatcher) CronFunc(cronExpr *CronExpr, _cb func()) *Cron {
	c := new(Cron)

	now := time.Now()
	nextTime := cronExpr.Next(now)
	if nextTime.IsZero() {
		return c
	}

	// callback
	var cb func()
	cb = func() {
		defer _cb()

		now := time.Now()
		nextTime := cronExpr.Next(now)
		if nextTime.IsZero() {
			return
		}
		c.t = disp.AfterFunc(nextTime.Sub(now), cb)
	}

	c.t = disp.AfterFunc(nextTime.Sub(now), cb)
	return c
}
