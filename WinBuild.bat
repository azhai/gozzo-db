@ECHO OFF

del dbtool.exe

go.exe build -ldflags="-s -w" -o dbtool.exe .

PAUSE