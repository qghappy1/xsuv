
package waitGroup

import(
	"sync"
    "os"
    "os/signal"
    "syscall"
	"time"
	"xsuv/util/log"
)

var id = uint(0)

var wait = &waitGroupT {
	ch:	make(chan bool), 
	wg:	&sync.WaitGroup{},
}

func init() {
	//logFile, err := os.OpenFile("crash1.log", os.O_CREATE|os.O_APPEND|os.O_RDWR, 0660)
	//if err != nil {
	//	log.Fatal("打开异常日志文件失败:%v", err)
	//    return
	//}
	// 将进程标准出错重定向至文件，进程崩溃时运行时将向该文件记录协程调用栈信息
	//syscall.Dup2(int(logFile.Fd()), int(os.Stderr.Fd())) 
	
    ch := make(chan os.Signal, 1)
    go func(){
	    <-ch
	    wait.ch <- true
	    log.Debug("recv signal stop server")
    }()
    signal.Notify(ch, syscall.SIGINT, syscall.SIGTERM)
}

type waitGroupT struct {
    ch        chan bool
    wg *sync.WaitGroup	
}

func IsSigStop() bool {
	select {
		case <-wait.ch:
			return true 
		default:		
	}
	return false 
}

func SigStop(){
	close(wait.ch)
	wait.wg.Wait()
}

func SigWait(){
    for {
		if !IsSigStop() {
			time.Sleep(time.Second*1)
		}else{
			break
		}
	}
}

func GoWrap(cb func()) {
	id = id + 1
	nid := id
	wait.wg.Add(1)
	log.DebugDepth(2, "start goroutine id:%v", nid)
	go func(){
		defer func() {
			if r := recover(); r != nil {
				log.Error("invoke goroutine id:%v error:%v", nid, r)
				log.ErrorStack()
				time.Sleep(time.Millisecond * 100) // 让日志可以完全输出
			}
		}()		
		cb()
		wait.wg.Done()
		log.Debug("end goroutine id:%v", nid)
	}()
}

func GoWrapEx(cb func(), fi func()) {
	id = id + 1
	nid := id
	wait.wg.Add(1)
	log.DebugDepth(2, "start goroutine id:%v", nid)
	go func(){
		defer func() {
			if r := recover(); r != nil {
				log.Error("invoke goroutine:%v error:%v", nid, r)
				log.ErrorStack()
				fi()
				time.Sleep(time.Millisecond * 100) // 让日志可以完全输出
			}
		}()		
		cb()
		wait.wg.Done()
		log.Debug("end goroutine id:%v", nid)
	}()
}

