
package conf

import (
	//"fmt"
	"bufio"
	"io"
	"os"
	"strings"
	//"errors"
	"path/filepath"
	//"flex/log"
)

type Conf struct {
	kvs map[string]string
	nodes map[string][]*Conf
}

func (this *Conf) GetValue(key string) string {
	if v, ok := this.kvs[key]; ok {
		return v
	}else{
		return ""
	}
}

func (this *Conf) GetConfs(key string) []*Conf {
	if v, ok := this.nodes[key]; ok {
		return v 
	}else{
		return nil 
	}
}

func (this *Conf) GetConf(key string) *Conf {
	if v, ok := this.nodes[key]; ok {
		if len(v) == 0 {
			return nil 
		}else{
			return v[0]
		}
	}else{	
		return nil 
	}
}

func Read(filename string) (*Conf, error) {
	ls, err := readFile(filename)
	if err!=nil {
		return nil, err 
	}
	conf, err := analysis(ls)
	return conf, err 
}

func newConf() *Conf {
	return &Conf{
		kvs : make(map[string]string),
		nodes : make(map[string][]*Conf),
	}
}

func getPath(fulleFilename string) string {
	fulleFilename, _ = filepath.Abs(fulleFilename)
	i := strings.LastIndex(fulleFilename, "\\")
	if i == -1 {
		i = strings.LastIndex(fulleFilename, "/")
	}
	return fulleFilename[:i+1] 
}

func readFile(filename string) ([]string, error) {
	f, err := os.Open(filename)
	defer f.Close()
	if err!=nil {
		return nil, err
	}
	filepath := getPath(filename)
	reader := bufio.NewReader(f)
	lines := make([]string, 0)
	for{
		line, err := reader.ReadString('\n')
		if line != "" {
			line = strings.Trim(line, " ")
			line = strings.Trim(line, "\t")
			line = strings.Trim(line, "\n")
			line = strings.Trim(line, "\r")
			if len(line) == 0 || line[0] == ';' || line[0] == byte(13) {
				continue
			}
			if strings.Contains(line, "#include") {
				line = strings.Trim(line, "#include")
				line = strings.Trim(line, " ")
				line = strings.Trim(line, "\t")
				line = strings.Trim(line, "\"")
				//fmt.Println(filepath+line)
				partlines, err := readFile(filepath+line)
				if err != nil {
					return nil, err 
				}
				lines = append(lines, partlines...)
			} else{
				bjmp := false 
				for i, c := range line {
					if c == ';' {
						lines = append(lines, line[0:i])
						bjmp = true 
						break
					}
				}
				if !bjmp {
					lines = append(lines, line)
				}	
			}
		}
		if err == io.EOF {
			break
		}
	}
	return lines, nil 
}

func analysis(lines []string) (*Conf, error) {
	conf := newConf()
	scan := newScanner(conf)
	for i:=0; i<len(lines); i++{
		//fmt.Println(lines[i])
		if err := scan.step([]byte(lines[i]), scan); err!=nil {
			return nil, err 
		}
	}	
	return conf, nil  
}






