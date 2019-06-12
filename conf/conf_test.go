package conf

import (
	"fmt"
	"testing"
)


// go test -v xsuv\conf
func Test_Conf(t *testing.T){
	root, _ := Read("E:/go/src/flex/conf/login.info")
	//_, err := Read("E:/go/src/flex/conf/login.info")
	
	fmt.Println(root.GetValue("md5key"))
	accounts := root.GetConf("accounts_db")
	fmt.Printf("db:%s\n", accounts.GetValue("connection_string"))
	conf2 := root.GetConf("login")
	if conf2 != nil {
		confs3 := conf2.GetConf("listen")
		if confs3 != nil {
			//fmt.Println(confs3.GetConf(""))
			fmt.Printf("address:%s\n", confs3.GetConf("").GetValue("port"))
			fmt.Printf("address:%s\n", confs3.GetConfs("")[1].GetValue("port"))
		}
	}
}
