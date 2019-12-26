@ECHO OFF

del gen2model.exe sync2table.exe dump2file.exe

go.exe build -ldflags="-s -w" ./cmd/gen2model/
gen2model.exe -f settings.toml
go.exe build -ldflags="-s -w" ./cmd/sync2table/
go.exe build -ldflags="-s -w" ./cmd/dump2file/

PAUSE