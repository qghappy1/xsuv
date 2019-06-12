SET API=.\api
::apix为服务器特有消息，与客户端不共用
SET APIX=.\apix
SET APIOUT=..\api
SET BIN=.\bin

echo build server api ...
del %APIOUT%\*.pb.go
%BIN%\protoc.exe -I%API% --go_out=%APIOUT% %API%\*.proto
::%BIN%\protoc.exe -I%APIX% --go_out=%APIOUT% %APIX%\*.proto
go install %APIOUT%

pause



