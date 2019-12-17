module github.com/azhai/gozzo-db

go 1.13

replace (
	golang.org/x/net => github.com/golang/net v0.0.0-20191125084936-ffdde1057850
	golang.org/x/sync => github.com/golang/sync v0.0.0-20190911185100-cd5d95a43a6e
	golang.org/x/tools => github.com/golang/tools v0.0.0-20191125224844-73cd2cc3b550
	golang.org/x/xerrors => github.com/golang/xerrors v0.0.0-20191011141410-1b5146add898
)

require (
	github.com/BurntSushi/toml v0.3.1
	github.com/codemodus/kace v0.5.1
	github.com/go-errors/errors v1.0.1
	github.com/go-sql-driver/mysql v1.4.1
	github.com/jinzhu/gorm v1.9.11
	github.com/jinzhu/inflection v1.0.0
	github.com/kr/pretty v0.1.0 // indirect
	github.com/lib/pq v1.2.0 // indirect
	github.com/stretchr/objx v0.2.0 // indirect
	github.com/stretchr/testify v1.4.0
	golang.org/x/net v0.0.0-20191209160850-c0dbc17a3553 // indirect
	golang.org/x/sync v0.0.0-20190911185100-cd5d95a43a6e // indirect
	golang.org/x/tools v0.0.0-20191216052735-49a3e744a425
	golang.org/x/xerrors v0.0.0-20191204190536-9bdfabe68543 // indirect
	google.golang.org/appengine v1.6.5 // indirect
	gopkg.in/yaml.v2 v2.2.7 // indirect
)
