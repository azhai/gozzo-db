@ECHO OFF

del code2mysql.exe table2file.exe

go.exe build -ldflags="-s -w" ./cmd/code2mysql/
go.exe build -ldflags="-s -w" ./cmd/table2file/

PAUSE