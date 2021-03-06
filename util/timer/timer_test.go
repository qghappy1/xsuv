package timer

import (
	"sync"
	"testing"
	"time"

	"github.com/qghappy1/xsuv/util/log"
)

var (
	once sync.Once
)

func testTimer() {
	Start()
	Tick(1, func() {
		go log.Debug("1")
	})
}

func Test_Timer(t *testing.T) {
	once.Do(testTimer)
	for {
		time.Sleep(1)
	}
}
