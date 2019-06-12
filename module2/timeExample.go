package module2

import (
	"fmt"
	"time"
)

func exampleTimer() {
	d := newDispatcher(10)

	// timer 1
	d.AfterFunc(1, func() {
		fmt.Println("My name is Leaf")
	})

	// timer 2
	t := d.AfterFunc(1, func() {
		fmt.Println("will not print")
	})
	t.Stop()

	// dispatch
	(<-d.ChanTimer).Cb()

	// Output:
	// My name is Leaf
}

func exampleCronExpr() {
	cronExpr, err := NewCronExpr("0 * * * *")
	if err != nil {
		return
	}

	fmt.Println(cronExpr.Next(time.Date(
		2000, 1, 1,
		20, 10, 5,
		0, time.UTC,
	)))

	// Output:
	// 2000-01-01 21:00:00 +0000 UTC
}

func exampleCron() {
	d := newDispatcher(10)

	// cron expr
	//cronExpr, err := NewCronExpr("*/2 * * * * *")
	// 秒, 分，时，日，月， 星期
	// 从2分钟开始之后每一分钟都执行该动作
	// *表示匹配所有值，/表示增长间隔
	cronExpr, err := NewCronExpr("0 2/1 * * * *")
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("Begin", time.Now().Unix())
	// cron
	var c *Cron
	c = d.CronFunc(cronExpr, func() {
		fmt.Println("time:", time.Now().Unix())
	})
	t := time.After(90*time.Second)
	for {
		select {
		case c := <-d.ChanTimer:
			c.Cb()
		case <-t:
			fmt.Println("exist")
			c.Stop()
			fmt.Println("End", time.Now().Unix())
			return
		}
	}
}
