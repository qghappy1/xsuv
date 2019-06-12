package module2

import (
	"fmt"
	"time"
)


func exampleFrame(){
	p := NewFrame(1000, func(){
		fmt.Println(time.Now().Unix(), ":frame update")
	})
	for i:=0; i<5; i++{
		p.Post(func(){
			fmt.Println(time.Now().Unix(), ":hello world")
			panic(1)
		})
		time.Sleep(time.Second)
	}
	p.Close()
	fmt.Println(time.Now().Unix(), ":close")
	time.Sleep(3*time.Second)
}