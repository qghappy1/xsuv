
package conf

import (
	//"fmt"
	//"strconv"
	"container/list"
	"strings"
)

const (
	STAT_NULL = iota
	STAT_KEY_BEGIN
	STAT_KEY_END					// preState is KEY_BEGIN, cur char is ' ' or '\n'
	STAT_VALUE_BEGIN				// preState is KEY_END
	STAT_VALUE_END					// preState is VALUE_BEGIN or STRING_END, cur char is '\n'
	STAT_CONF_BEGIN					// preState is KEY_END, cur char is '{'
	STAT_CONF_END					// cur char is '}'
	STAT_NOTE_BEGIN					// cur char is ';'
	STAT_NOTE_END					// preState is NOTE_BEGIN, cur char is '\n'
	STAT_STRING_BEGIN				// preState is VALUE_BEGIN, cur char is '"'
	STAT_STRING_END					// preState is STRING_BEGIN, cur char is '"'
)

type scanner struct {
	curState int 
	preState int 
	curKey string
	curConf *Conf
	step func(data []byte, scan *scanner) error
	err error 
	stackConf *list.List
}

func newScanner(conf *Conf) *scanner {
	this := new(scanner)
	this.curState = STAT_KEY_BEGIN
	this.preState = STAT_NULL
	this.step = stepKey
	this.curKey = ""
	this.curConf = conf
	this.stackConf = list.New()
	this.stackConf.PushBack(conf)
	return this 
}

func stepKey(data []byte, scan *scanner) error {
	if data[0] == '}' {
		//fmt.Printf("}next stepKey\n")
		scan.preState = scan.curState
		scan.curState = STAT_CONF_END
		back := scan.stackConf.Back()
		scan.stackConf.Remove(back)
		scan.curConf = scan.stackConf.Back().Value.(*Conf)
		scan.step = stepKey
		return nil 		
	}
	if data[0] == ';' {
		return nil 
	}
	scan.preState = scan.curState
	scan.curState = STAT_VALUE_BEGIN
	scan.step = stepVaule
	for i, c := range data {
		if c == ' ' || c == '\t' {
			scan.curKey = string(data[:i])
			tmp := string(data[i+1:])
			tmp = strings.Trim(tmp, " ")
			tmp = strings.Trim(tmp, "\t")
			tmp = strings.Trim(tmp, "\"")
			tmpbytes := []byte(tmp)
			//fmt.Printf("' 'next stepVaule key:%s\n", scan.curKey)	
			return scan.step(tmpbytes, scan) 
		}
		if c == ';' {
			scan.curKey = string(data[:i])
			//fmt.Printf(";next stepVaule key:%s\n", scan.curKey)
			return nil
		}		
	}
	scan.curKey = string(data)
	scan.curKey = strings.Trim(scan.curKey, "\"")
	//fmt.Printf("stepKey key:%s\n", scan.curKey)
	return nil 	
}

func stepVaule(data []byte, scan *scanner) error {
	if len(data) == 0 {
		scan.curConf.kvs[scan.curKey] = ""
		scan.preState = STAT_VALUE_END
		scan.curState = STAT_KEY_BEGIN
		scan.step = stepKey
		return nil 			
	}
	if data[0] == '{' {
		//fmt.Printf("{next stepKey\n")
		scan.preState = scan.curState
		scan.curState = STAT_CONF_BEGIN
		conf := newConf() 
		//fmt.Printf("new conf:%s\n", scan.curKey)
		if arrConf, ok := scan.curConf.nodes[scan.curKey]; ok {
			scan.curConf.nodes[scan.curKey] = append(arrConf, conf)
		} else {
			arrConf := make([]*Conf, 1)
			arrConf[0] = conf
			scan.curConf.nodes[scan.curKey] = arrConf
		}
		scan.curConf = conf
		scan.stackConf.PushBack(conf)
		scan.step = stepKey
		return nil 
	}
	for i, c := range data {
		if c == ';' {
			scan.curConf.kvs[scan.curKey] = string(data[:i])
			scan.preState = STAT_VALUE_END
			scan.curState = STAT_KEY_BEGIN
			scan.step = stepKey
			return nil
		}		
	}
	//fmt.Printf("next stepKey k:%s v:%s\n", scan.curKey, string(data))
	scan.curConf.kvs[scan.curKey] = string(data)
	scan.preState = STAT_VALUE_END
	scan.curState = STAT_KEY_BEGIN
	scan.step = stepKey
	
	//if scan.curKey == "port" {
		//fmt.Printf("next stepKey k:%s v:%s\n", scan.curKey, string(data))
		//fmt.Println(scan.curConf)
	//}
	
	return nil 	
}




