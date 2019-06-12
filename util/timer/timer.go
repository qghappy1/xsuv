package timer

import (
	"fmt"
	"time"

	"github.com/qghappy1/xsuv/util/waitGroup"
)

const formatTime = "2006-01-02 15:04:05"

// 1s定时器
var timer = createTimer(time.Millisecond * 1000)

// 节点
type Node struct {
	tm int
	f  func()
}

type Timer struct {
	events []*Node
	tick   time.Duration
	quit   bool
}

func (n *Node) String() string {
	return fmt.Sprintf("Node.time,%d", n.tm)
}

// 创建定时器
func createTimer(d time.Duration) *Timer {
	t := new(Timer)
	t.tick = d
	t.quit = false
	t.events = make([]*Node, 0)
	return t
}

func (t *Timer) NewTimer(tm int, f func()) *Node {
	if t.quit == false {
		n := new(Node)
		n.f = f
		n.tm = tm
		//tt := time.Unix(int64(tm), 0)
		//timeNow := tt.Format(formatTime)
		//log.Debug("add timer:%s", timeNow)
		t.events = append(t.events, n)
		return n
	} else {
		return nil
	}
}

func (t *Timer) update() {
	now := int(time.Now().Unix())
	i := 0
	for {
		if i >= len(t.events) {
			break
		}
		e := t.events[i]
		//log.Debug("own timer:%d", e.tm)
		if e.tm <= now {
			e.f()
			//log.Debug("process timer:%d", e.tm)
			t.events = append(t.events[:i], t.events[i+1:]...)
		} else {
			i++
		}
	}
}

func (t *Timer) Start() {
	tick := time.NewTicker(t.tick)
	waitGroup.GoWrap(func() {
		defer tick.Stop()
		for {
			time.Sleep(t.tick)
			t.update()
			if t.quit {
				break
			}
		}
	})
}

func (t *Timer) Stop() {
	t.quit = true
}

func Start() {
	timer.Start()
}

func AfterFunc(second int, f func()) {
	//now, _ := time.ParseInLocation(formatTime, time.Now().Format(format), time.Local)
	//log.Debug("after process:%s", now)
	timer.NewTimer(int(time.Now().Unix())+second, f)
}

func Tick(second int, f func()) {
	timer.NewTimer(int(time.Now().Unix())+second, func() {
		f()
		Tick(second, f)
	})
}

// 每天定时处理
func TimeHourFunc(hour int, f func()) {
	tm := time.Now().Unix()
	tnext := 0
	t := time.Unix(tm, 0)

	// 6点前启动则到6点执行，6点后启动则第二天6点执行 hour = 6
	if t.Hour() > hour {
		t = time.Unix(tm+24*60*60, 0)
	}
	tnext = int(time.Date(t.Year(), t.Month(), t.Day(), hour, 0, 0, 0, time.Local).Unix())

	//tt := time.Unix(int64(tnext), 0)
	//timeNow := tt.Format(formatTime)
	//log.Debug("timer hour:%s", timeNow)
	timer.NewTimer(tnext, func() {
		f()
		Tick(24*60*60, f)
	})
}

// 整点的某分钟处理
func TimeMinuteFunc(minute int, f func()) {
	//now, _ := time.ParseInLocation(formatTime, time.Now().Format(format), time.Local)
	//log.Debug("time minute process:%s", now1)

	now, _ := time.ParseInLocation(formatTime, time.Now().Format(formatTime), time.Local)
	t2 := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.Local)

	for i := 0; i < 24; i++ {
		t2 = time.Date(now.Year(), now.Month(), now.Day(), i, minute, 0, 0, time.Local)
		if t2.Hour() < now.Hour() {

			continue
		} else if t2.Hour() == now.Hour() {
			if t2.Before(now) {
				t2 = time.Date(now.Year(), now.Month(), now.Day(), i+1, minute, 0, 0, time.Local)
				d := t2.Sub(now)
				//fmt.Println("time:", t2)
				//fmt.Println("now:", now)
				//fmt.Println("d:", d)
				timer.NewTimer(int(time.Now().Unix())+int(d.Seconds()), func() {
					f()
					Tick(60*minute, f)
				})
			} else {
				d := t2.Sub(now)
				//fmt.Println("time2:", t2)
				//fmt.Println("now:", now)
				//fmt.Println("d:", d)
				timer.NewTimer(int(time.Now().Unix())+int(d.Seconds()), func() {
					f()
					Tick(60*minute, f)
				})
			}
			break
		} else {

		}
	}
}
