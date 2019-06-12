
set GOOS=windows
set GOPACH=amd64

go build -o bin/gnatsd.exe main.go
::go build  -gcflags "-N -l"

pause