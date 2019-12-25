@ECHO OFF

del table2file.exe load4toml.exe

go.exe build -ldflags="-s -w" ./cmd/table2file/
table2file.exe -f settings.toml
go.exe build -ldflags="-s -w" ./cmd/load4toml/

PAUSE