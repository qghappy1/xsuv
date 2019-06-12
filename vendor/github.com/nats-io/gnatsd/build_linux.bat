


set GOOS=linux
set GOPACH=amd64


go build -o bin/gnatsd main.go
::go build  -gcflags "-N -l"

pause