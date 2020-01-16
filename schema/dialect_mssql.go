package schema

import (
	"fmt"

	"github.com/azhai/gozzo-utils/common"
	"github.com/azhai/gozzo-utils/redisw"
)

const MSSQL_DEFAULT_PORT uint16 = 1433

type Mssql struct {
}

func (Mssql) GetDSN(params redisw.ConnParams) (string, string) {
	dsn := "sqlserver://"
	user := redisw.ConcatWith(params.Username, params.Password)
	if user != "" {
		dsn += user + "@"
	}
	dsn += params.GetAddr("127.0.0.1", MSSQL_DEFAULT_PORT)
	if params.Database != "" {
		dsn += "?database" + params.Database
	}
	return "mssql", dsn
}

func (Mssql) QuoteIdent(ident string) string {
	return common.WrapWith(ident, "[", "]")
}

func (Mssql) dbNameCol() string {
	return "DB_NAME()"
}

func (d Mssql) dbNameVal(dbname string) string {
	if dbname == "" {
		return d.dbNameCol()
	} else {
		return common.WrapWith(dbname, "'", "'")
	}
}

func (d Mssql) CurrDbNameSql() string {
	return fmt.Sprintf("SELECT %s", d.dbNameCol())
}

func (Mssql) tableNameTpl() string {
	return common.ReduceSpaces(`
			SELECT table_name, table_catalog, table_rows, table_comment
			FROM
				information_schema.tables
			WHERE
				table_type = '%s' AND table_catalog %s
		`)
}

func (d Mssql) TableNameSql(dbname string, more bool) string {
	if more {
		dbcond := "LIKE " + common.WrapWith(dbname, "'", "%'")
		return fmt.Sprintf(d.tableNameTpl(), "BASE TABLE", dbcond)
	} else {
		dbcond := "= " + d.dbNameVal(dbname)
		return fmt.Sprintf(d.tableNameTpl(), "BASE TABLE", dbcond)
	}
}

func (d Mssql) ViewNameSql(dbname string, more bool) string {
	if more {
		dbcond := "LIKE " + common.WrapWith(dbname, "'", "%'")
		return fmt.Sprintf(d.tableNameTpl(), "VIEW", dbcond)
	} else {
		dbcond := "= " + d.dbNameVal(dbname)
		return fmt.Sprintf(d.tableNameTpl(), "VIEW", dbcond)
	}
}

func (Mssql) ColumnTypeSql(fullTableName string) string {
	return fmt.Sprintf("SELECT * FROM %s WHERE 1=0", fullTableName)
}

func (d Mssql) ColumnInfoSql(table, dbname string) string {
	tpl := common.ReduceSpaces(`
			SELECT 
				column_name, column_type, column_key,
				column_default, extra, column_comment
			FROM
				information_schema.columns
			WHERE
				table_name = '%s' AND table_schema = %s
			ORDER BY ordinal_position
		`)
	dbname = d.dbNameVal(dbname)
	return fmt.Sprintf(tpl, table, dbname)
}
