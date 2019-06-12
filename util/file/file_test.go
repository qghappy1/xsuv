package file

import (
	"testing"
	"fmt"
)

func Test_file(t *testing.T){
	fmt.Println(GetCurFilePath())
	fmt.Println(IsFileExist("E:/go/src/flex/util/file.go"))
	//fmt.Println(RenameFile("E:/go/src/flex/util/file.go", "E:/go/src/flex/util/file2.go"))
}