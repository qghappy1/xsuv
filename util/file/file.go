
package file

import (
	"os"
	"path/filepath"
	"os/exec"
	"strings"
)

func GetPath(fulleFilename string) string {
	i := strings.LastIndex(fulleFilename, "\\")
	return fulleFilename[:i+1] 
}

func GetCurFilePath() string {
	file, _ := exec.LookPath(os.Args[0])
    fulleFilename, _ := filepath.Abs(file)
	i := strings.LastIndex(fulleFilename, "\\")
	return fulleFilename[:i+1] 
}

func GetFileSize(file *os.File) (int64, error) {
	finfo, err := file.Stat()
	if err != nil {
		return 0, err 
	}
	return finfo.Size(), nil 
}

func IsFileExist(filename string) bool {
	f, err := os.Open(filename)
	defer f.Close()
	if err!=nil && os.IsNotExist(err){
		return false 				
	}
	return true 
}

func RenameFile(oldfilename string, newfilename string) bool {
	f, err := os.Open(newfilename)
	defer f.Close()
	if err!=nil && os.IsNotExist(err){
		err = os.Rename(oldfilename, newfilename)
		if err != nil {
			return true 
		}
	}
	return false 
}