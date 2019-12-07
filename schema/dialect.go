package schema

import (
	"database/sql"
	"fmt"
	"strings"
)

var drivers = map[string]string{
	"*mysql.MySQLDriver":          "mysql",    // github.com/go-sql-driver/mysql
	"*godrv.Driver":               "mysql",    // github.com/ziutek/mymysql - TODO(js) No datatypes.
	"*pq.Driver":                  "postgres", // github.com/lib/pq
	"*stdlib.Driver":              "postgres", // github.com/jackc/pgx
	"*pgsqldriver.postgresDriver": "postgres", // github.com/jbarham/gopgsqldriver - TODO(js) No datatypes.
	"*sqlite3.SQLiteDriver":       "sqlite",   // github.com/mattn/go-sqlite3
	"*sqlite.impl":                "sqlite",   // github.com/gwenn/gosqlite
	"sqlite3.Driver":              "sqlite",   // github.com/mxk/go-sqlite - TODO(js) No datatypes.
	"*mssql.Driver":               "mssql",    // github.com/denisenkom/go-mssqldb
	"*mssql.MssqlDriver":          "mssql",    // github.com/denisenkom/go-mssqldb
	"*freetds.MssqlDriver":        "mssql",    // github.com/minus5/gofreetds - TODO(js) No datatypes. Error on create view.
	"*goracle.drv":                "oracle",   // gopkg.in/goracle.v2
	"*ora.Drv":                    "oracle",   // gopkg.in/rana/ora.v4 - TODO(js) Mismatched datatypes.
	"*oci8.OCI8DriverStruct":      "oracle",   // github.com/mattn/go-oci8 - TODO(js) Mismatched datatypes.
	"*oci8.OCI8Driver":            "oracle",   // github.com/mattn/go-oci8 - TODO(js) Mismatched datatypes.
}

var dialects = map[string]Dialect{
	"mssql":    &Mssql{},
	"mysql":    &Mysql{},
	"oracle":   &Oracle{},
	"postgres": &Postgres{},
	"sqlite":   &Sqlite{},
}

type FetchFunc = func(rows *sql.Rows) (err error)

type Dialect interface {
	GetDSN(params ConnParams) (string, string)
	QuoteIdent(ident string) string
	CurrDbNameSql() string
	TableNameSql(dbname string, more bool) string
	ViewNameSql(dbname string, more bool) string
	ColumnTypeSql(fullTableName string) string
	ColumnInfoSql(table, dbname string) string
}

func GetDialectByName(name string) Dialect {
	name = strings.ToLower(name)
	if d, ok := dialects[name]; ok {
		return d
	}
	return nil
}

func GetDialect(db *sql.DB) (Dialect, string, error) {
	drv := fmt.Sprintf("%T", db.Driver())
	if name, ok := drivers[drv]; ok {
		if d := GetDialectByName(name); d != nil {
			return d, name, nil
		}
	}
	err := UnknownDriverError{Driver: drv}
	return nil, "", err
}

type Schema struct {
	db         *sql.DB
	tables     map[string]map[string]TableInfo
	DriverName string
	Dialect    Dialect
	Error      error
}

func NewSchema(db *sql.DB) *Schema {
	s := &Schema{db: db, tables: make(map[string]map[string]TableInfo)}
	s.Dialect, s.DriverName, s.Error = GetDialect(db)
	return s
}

func (s *Schema) Query(fetch FetchFunc, dsql string, args ...interface{}) error {
	if s.Error != nil {
		return s.Error
	}
	var rows *sql.Rows
	rows, s.Error = s.db.Query(dsql, args...)
	defer func() {
		s.Error = rows.Close()
	}()
	if s.Error != nil {
		return s.Error
	}
	return fetch(rows)
}

func (s *Schema) GetStrings(dsql string) (names []string) {
	_ = s.Query(func(rows *sql.Rows) (err error) {
		for rows.Next() {
			var name string
			if err = rows.Scan(&name); err == nil {
				names = append(names, name)
			}
		}
		return
	}, dsql)
	return
}

func (s *Schema) GetString(dsql string, offset int) (name string) {
	if s.Error != nil {
		return
	}
	row := s.db.QueryRow(dsql)
	if offset == 1 {
		s.Error = row.Scan(nil, &name)
	} else if offset == 2 {
		s.Error = row.Scan(nil, nil, &name)
	} else if offset == 3 {
		s.Error = row.Scan(nil, nil, nil, &name)
	} else {
		s.Error = row.Scan(&name)
	}
	return
}

func (s *Schema) GetCurrDbName() (dbname string) {
	dsql := s.Dialect.CurrDbNameSql()
	if dsql == "" {
		return
	}
	if s.DriverName == "sqlite" {
		dbname = s.GetString(dsql, 1)
	} else {
		dbname = s.GetString(dsql, 0)
	}
	return
}

func (s *Schema) GetTableNames(dbname string) (names []string) {
	if s.DriverName == "mysql" || s.DriverName == "postgres" {
		if dbname == "" {
			dbname = s.GetCurrDbName()
		}
		if _, ok := s.tables[dbname]; !ok {
			s.tables[dbname] = s.AllTableInfos(dbname, false)
		}
		for name := range s.tables[dbname] {
			names = append(names, name)
		}
	} else {
		dsql := s.Dialect.TableNameSql(dbname, false)
		names = s.GetStrings(dsql)
	}
	return
}

func (s *Schema) GetViewNames(dbname string) (names []string) {
	dsql := s.Dialect.ViewNameSql(dbname, false)
	names = s.GetStrings(dsql)
	return
}

func (s *Schema) AllTableInfos(dbname string, more bool) map[string]TableInfo {
	ti := TableInfo{Quote: s.Dialect.QuoteIdent}
	dsql := s.Dialect.TableNameSql(dbname, more)
	infos := make(map[string]TableInfo)
	_ = s.Query(func(rows *sql.Rows) (err error) {
		for rows.Next() {
			if s.DriverName == "sqlite" {
				err = rows.Scan(&ti.TableName)
			} else if s.DriverName == "oracle" {
				err = rows.Scan(&ti.TableName, &ti.DbName, &ti.TableRows)
			} else {
				err = rows.Scan(&ti.TableName, &ti.DbName, &ti.TableRows, &ti.TableComment)
			}
			if ti.TableName != "" {
				infos[ti.TableName] = ti
			}
		}
		return
	}, dsql)
	return infos
}

func (s *Schema) ListTable(dbname string, refresh bool) (tables map[string]TableInfo) {
	var ok bool
	if dbname == "" {
		dbname = s.GetCurrDbName()
	}
	more := strings.HasSuffix(dbname, "%")
	if tables, ok = s.tables[dbname]; !ok || more || refresh {
		tables = s.AllTableInfos(dbname, more)
		for name, info := range tables {
			if _, ok = s.tables[info.DbName]; !ok {
				s.tables[info.DbName] = make(map[string]TableInfo)
			}
			s.tables[info.DbName][name] = tables[name]
		}
	}
	return
}

func (s *Schema) GetTableInfo(table, dbname string) (tbInfo TableInfo) {
	var ok bool
	if dbname == "" {
		dbname = s.GetCurrDbName()
	}
	if _, ok = s.tables[dbname]; !ok {
		s.tables[dbname] = s.AllTableInfos(dbname, false)
	}
	if tbInfo, ok = s.tables[dbname][table]; !ok {
		tbInfo = TableInfo{
			DbName: dbname,
			Quote:  s.Dialect.QuoteIdent,
		}
	}
	return
}

func (s *Schema) GetColumnTypes(fullTableName string) (cols []*sql.ColumnType) {
	if s.Error != nil {
		return
	}
	dsql := s.Dialect.ColumnTypeSql(fullTableName)
	var rows *sql.Rows
	rows, s.Error = s.db.Query(dsql)
	defer rows.Close()
	if s.Error != nil {
		return
	}
	cols, s.Error = rows.ColumnTypes()
	return
}

func (s *Schema) GetColumnExtras(table, dbname string) map[string]ColumnExtra {
	var rows *sql.Rows
	dsql := s.Dialect.ColumnInfoSql(table, dbname)
	rows, s.Error = s.db.Query(dsql)
	defer rows.Close()
	if s.Error != nil {
		return nil
	}
	extras := make(map[string]ColumnExtra)
	for rows.Next() {
		ce := DefaultColumnExtra
		_ = rows.Scan(&ce.FieldName, &ce.FullType, &ce.IndexType, &ce.DefaultValue,
			&ce.Extra, &ce.Comment, &ce.MaxSize, &ce.PrecSize, &ce.PrecScale)
		extras[ce.FieldName] = ce
	}
	return extras
}

func (s *Schema) GetColumnInfos(table, dbname string) []*ColumnInfo {
	if dbname == "" {
		dbname = s.GetCurrDbName()
	}
	tbInfo := s.GetTableInfo(table, dbname)
	if tbInfo.TableName == "" {
		tbInfo.TableName = table
	}
	extras := s.GetColumnExtras(table, dbname)
	fullTableName := tbInfo.GetFullName(true)
	cols := s.GetColumnTypes(fullTableName)
	if s.Error != nil {
		return nil
	}
	infos := make([]*ColumnInfo, len(cols))
	for i, ct := range cols {
		infos[i] = &ColumnInfo{
			Table:       &tbInfo,
			ColumnType:  ct,
			ColumnExtra: extras[ct.Name()],
		}
	}
	return infos
}
