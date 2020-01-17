module github.com/azhai/gozzo-db

go 1.13

replace (
	golang.org/x/net => github.com/golang/net v0.0.0-20200114155413-6afb5195e5aa
	golang.org/x/sync => github.com/golang/sync v0.0.0-20190911185100-cd5d95a43a6e
	golang.org/x/tools => github.com/golang/tools v0.0.0-20200116062425-473961ec044c
	golang.org/x/xerrors => github.com/golang/xerrors v0.0.0-20191011141410-1b5146add898
)

require (
	github.com/BurntSushi/toml v0.3.1
	github.com/azhai/gozzo-utils v0.4.2
	github.com/codemodus/kace v0.5.1
	github.com/go-errors/errors v1.0.1
	github.com/go-sql-driver/mysql v1.5.0 // indirect
	github.com/gomodule/redigo v2.0.0+incompatible
	github.com/jinzhu/gorm v1.9.12
	github.com/jinzhu/inflection v1.0.0
	github.com/stretchr/testify v1.4.0
	golang.org/x/tools v0.0.0-20200116062425-473961ec044c
)
