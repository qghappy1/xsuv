package lfs

import (
	"os"
	"syscall"

	"github.com/yuin/gopher-lua"
)

func attributesFill(tbl *lua.LTable, stat os.FileInfo) {
	sys, ok := stat.Sys().(*syscall.Win32FileAttributeData)
	if !ok {
		return
	}
	{
		var mode string
		if stat.IsDir() {
			mode = "directory"
		}else{
			mode = "file"
		}
		tbl.RawSetH(lua.LString("mode"), lua.LString(mode))
	}
	tbl.RawSetH(lua.LString("access"), lua.LNumber(sys.LastAccessTime.Nanoseconds()/1e9))
	tbl.RawSetH(lua.LString("modification"), lua.LNumber(sys.LastWriteTime.Nanoseconds()/1e9))
	tbl.RawSetH(lua.LString("change"), lua.LNumber(sys.CreationTime.Nanoseconds()/1e9))
	tbl.RawSetH(lua.LString("size"), lua.LNumber(stat.Size()))
}
