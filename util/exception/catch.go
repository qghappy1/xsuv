
package exception

import (
	"os"
	"xsuv/util/log"
)

func CatchError(f func()){
	if r := recover(); r != nil {
		log.ErrorStack()
		if f!=nil {
			f()
		}		
	}
}

func CatchAndExit(f func()){
	if r := recover(); r != nil {
		log.ErrorStack()
		if f!=nil {
			f()
		}		
		os.Exit(0)
	}
}
