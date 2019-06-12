package etcdv3

import (
	"fmt"
	"time"

	"github.com/qghappy1/xsuv/util/log"
)

func CheckVersion(c *Client, serviceName, version string, closeSig chan bool) {
	go func() {
		defer log.ErrorPanic()
		for {
			ver := Get(c, fmt.Sprintf("%vVersion", serviceName), 10)
			if ver == "" {
				Put(c, fmt.Sprintf("%vVersion", serviceName), version)
			} else {
				if ver > version {
					closeSig <- true
					break
				}
			}
			time.Sleep(15 * time.Second)
		}
	}()
}
