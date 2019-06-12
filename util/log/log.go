package log

import (
	"errors"
	"fmt"
	"log"
	"os"
	"runtime"
	"runtime/debug"
	"strings"
	"time"

	"github.com/qghappy1/xsuv/util/file"
)

const (
	DEBUG = 0
	INFO  = 1
	WARN  = 2
	ERROR = 3
)

var (
	filename = "debug"
	filepath = ""
	flag     = -1
)

func openLogFile() (*os.File, error) {
	now := time.Now()
	//exist := false
	//fname := ""
	for i := 0; i < 200000; i++ {
		fileName := fmt.Sprintf("%s%d%02d%02d_%d.log", filename, now.Year(), now.Month(), now.Day(), i)
		file, err := os.OpenFile(fileName, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0777)
		if err != nil {
			return nil, err
		}
		fileInfo, err := file.Stat()
		if err != nil {
			return nil, err
		}
		if fileInfo.Size() >= 20*1024*1024 {
			//exist = true
			//fname = fileName
			continue
		}
		log.SetOutput(file)
		//if exist {
		//	log.Println(fmt.Sprintf("pre fname:%v, cur:%v %d%02d%02d_%d tm:%v", fname, fileName, now.Year(), now.Month(), now.Day(), i, now.Unix() ))
		//}
		return file, err
	}
	return nil, errors.New("can not open log file")
}

func SetFilename(name string) {
	filepath = file.GetCurFilePath() + "log/"
	filename = file.GetCurFilePath() + "log/" + name
	os.MkdirAll(filepath, 0777)
}

func SetFlag(i int) {
	flag = i
}

func GetFlag() int {
	return flag
}

func Info(format string, v ...interface{}) {
	if flag > INFO {
		return
	}
	format = fmt.Sprintf("[INFO ] %s", format)
	s := fmt.Sprintf(format, v...)
	fmt.Println(s)
	file, err := openLogFile()
	if err != nil {
		return
	}
	defer file.Close()
	log.Println(s)
}

func LuaInfo(s string) {
	if flag > INFO {
		return
	}
	fmt.Println(s)
	file, err := openLogFile()
	if err != nil {
		return
	}
	defer file.Close()
	log.Println(s)
}

func Debug(format string, v ...interface{}) {
	if flag > DEBUG {
		return
	}
	_, file, line, ok := runtime.Caller(1)
	if !ok {
		file = "???"
		line = 0
	}
	i := strings.LastIndex(file, "/")
	file = file[i+1:]
	format = fmt.Sprintf("%s.%d[DEBUG] %s", file, line, format)
	s := fmt.Sprintf(format, v...)
	fmt.Println(s)
	f, err := openLogFile()
	if err != nil {
		return
	}
	defer f.Close()
	log.Println(s)
}

func LuaDebug(s string) {
	if flag > DEBUG {
		return
	}
	fmt.Println(s)
	f, err := openLogFile()
	if err != nil {
		return
	}
	defer f.Close()
	log.Println(s)
}

func DebugDepth(depth int, format string, v ...interface{}) {
	if flag > DEBUG {
		return
	}
	_, file, line, ok := runtime.Caller(depth)
	if !ok {
		file = "???"
		line = 0
	}
	i := strings.LastIndex(file, "/")
	file = file[i+1:]
	format = fmt.Sprintf("%s.%d[DEBUG] %s", file, line, format)
	s := fmt.Sprintf(format, v...)
	fmt.Println(s)
	f, err := openLogFile()
	if err != nil {
		return
	}
	defer f.Close()
	log.Println(s)
}

func ErrorDepth(depth int, format string, v ...interface{}) {
	_, file, line, ok := runtime.Caller(depth)
	if !ok {
		file = "???"
		line = 0
	}
	i := strings.LastIndex(file, "/")
	file = file[i+1:]
	format = fmt.Sprintf("%s.%d[ERROR] %s", file, line, format)
	s := fmt.Sprintf(format, v...)
	fmt.Println(s)
	f, err := openLogFile()
	if err != nil {
		return
	}
	defer f.Close()
	log.Println(s)
}

func FatalDepth(depth int, format string, v ...interface{}) {
	_, file, line, ok := runtime.Caller(depth)
	if !ok {
		file = "???"
		line = 0
	}
	i := strings.LastIndex(file, "/")
	file = file[i+1:]
	format = fmt.Sprintf("%s.%d[ERROR] %s", file, line, format)
	s := fmt.Sprintf(format, v...)
	fmt.Println(s)
	f, err := openLogFile()
	if err != nil {
		return
	}
	defer f.Close()
	log.Println(s)
	os.Exit(-1)
}

func WarnDepth(depth int, format string, v ...interface{}) {
	if flag > WARN {
		return
	}
	_, file, line, ok := runtime.Caller(depth)
	if !ok {
		file = "???"
		line = 0
	}
	i := strings.LastIndex(file, "/")
	file = file[i+1:]
	format = fmt.Sprintf("%s.%d[WARN] %s", file, line, format)
	s := fmt.Sprintf(format, v...)
	fmt.Println(s)
	f, err := openLogFile()
	if err != nil {
		return
	}
	defer f.Close()
	log.Println(s)
}

func Warn(format string, v ...interface{}) {
	if flag > WARN {
		return
	}
	_, file, line, ok := runtime.Caller(1)
	if !ok {
		file = "???"
		line = 0
	}
	i := strings.LastIndex(file, "/")
	file = file[i+1:]
	format = fmt.Sprintf("%s.%d[WARN] %s", file, line, format)
	s := fmt.Sprintf(format, v...)
	fmt.Println(s)
	f, err := openLogFile()
	if err != nil {
		return
	}
	defer f.Close()
	log.Println(s)
}

func Error(format string, v ...interface{}) {
	_, file, line, ok := runtime.Caller(1)
	if !ok {
		file = "???"
		line = 0
	}
	i := strings.LastIndex(file, "/")
	file = file[i+1:]
	format = fmt.Sprintf("%s.%d[ERROR] %s", file, line, format)
	s := fmt.Sprintf(format, v...)
	fmt.Println(s)
	f, err := openLogFile()
	if err != nil {
		return
	}
	defer f.Close()
	log.Println(s)
}

func LuaError(s string) {
	fmt.Println(s)
	f, err := openLogFile()
	if err != nil {
		return
	}
	defer f.Close()
	log.Println(s)
}

func Fatal(format string, v ...interface{}) {
	s := fmt.Sprintf(format, v...)
	fmt.Println(s)
	f, err := openLogFile()
	if err != nil {
		return
	}
	defer f.Close()
	log.Println(s)
	os.Exit(-1)
}

func ErrorStack() {
	Error("%s", string(debug.Stack()))
}

func FatalStack() {
	Fatal("%s", string(debug.Stack()))
}

func FatalPanic() {
	if r := recover(); r != nil {
		FatalStack()
		time.Sleep(time.Millisecond * 100) // 让日志可以完全输出
	}
}

func ErrorPanic() {
	if r := recover(); r != nil {
		ErrorStack()
		time.Sleep(time.Millisecond * 100) // 让日志可以完全输出
	}
}
