package module2

import (
	"os"
	"os/signal"
	"sync"

	"github.com/qghappy1/xsuv/util/log"
)

type modules struct {
	mi       IModule
	closeSig chan bool
	wg       sync.WaitGroup
}

var mods []*modules

func registerModule(mi IModule) {
	m := new(modules)
	m.mi = mi
	m.closeSig = make(chan bool, 1)

	mods = append(mods, m)
}

func initModules() {
	//log.Info("Server init")
	for i := 0; i < len(mods); i++ {
		mods[i].mi.OnInit()
	}
	//log.Info("Server run")
	for i := 0; i < len(mods); i++ {
		m := mods[i]
		m.wg.Add(1)
		go run(m)
	}
}

func destroyModules() {
	for i := len(mods) - 1; i >= 0; i-- {
		m := mods[i]
		m.closeSig <- true
		m.wg.Wait()
		destroy(m)
	}
}

func run(m *modules) {
	defer func() {
		if r := recover(); r != nil {
			log.FatalStack()
		}
	}()
	m.mi.Run(m.closeSig)
	m.wg.Done()
}

func destroy(m *modules) {
	defer func() {
		if r := recover(); r != nil {
			log.FatalStack()
		}
	}()

	m.mi.OnDestroy()
}

func Run(mods ...IModule) {
	log.Info("Server starting")
	// module
	for i := 0; i < len(mods); i++ {
		registerModule(mods[i])
	}
	initModules()
	// close
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, os.Kill)
	sig := <-c
	log.Info("Server closing down (signal: %v)", sig)
	destroyModules()
}
