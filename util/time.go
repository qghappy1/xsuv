
package util

import (
	"fmt"
	"time"
	"math/rand"
)

var s = rand.NewSource(time.Now().UnixNano())

func TimeToString(t time.Time) string {
	return t.Format("2006-01-02 15:04:05.000 -0700")
}

// TimeAfter(time.Now(), time.Now(), 10) return false 
func TimeAfter(t1 time.Time, t2 time.Time, second int) bool {
	d, err := time.ParseDuration(fmt.Sprintf("%ds", second))
	fmt.Println("%v", d)	
	if err != nil {
		fmt.Printf(err.Error())
		return false 
	}
	t2.Add(d)  
	return t1.After(t2)
}


func Today() time.Time {
	t1 := time.Now()
	return time.Date(t1.Year(),t1.Month(),t1.Day(),0,0,0,0, time.Local)
}

func Intn(n int) int {
	return rand.New(s).Intn(n)
}

