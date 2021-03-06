package schema

import (
	"fmt"

	"github.com/azhai/gozzo-utils/common"
	"github.com/azhai/gozzo-utils/redisw"
)

const MYSQL_DEFAULT_PORT uint16 = 3306

type Mysql struct {
}

func (Mysql) GetDSN(params redisw.ConnParams) (string, string) {
	user := redisw.ConcatWith(params.Username, params.Password)
	addr := params.GetAddr("127.0.0.1", MYSQL_DEFAULT_PORT)
	dsn := user + "@"
	if addr != "" {
		dsn += "(" + addr + ")"
	}
	if params.Database != "" {
		dsn += "/" + params.Database
	}
	dsn += "?parseTime=true&loc=Local"
	if charset, ok := params.Options["charset"]; ok {
		dsn += fmt.Sprintf("&charset=%s", charset)
	}
	return "mysql", dsn
}

func (Mysql) QuoteIdent(ident string) string {
	return common.WrapWith(ident, "`", "`")
}

func (Mysql) dbNameCol() string {
	return "DATABASE()"
}

func (d Mysql) dbNameVal(dbname string) string {
	if dbname == "" {
		return d.dbNameCol()
	} else {
		return common.WrapWith(dbname, "'", "'")
	}
}

func (d Mysql) CurrDbNameSql() string {
	return fmt.Sprintf("SELECT %s", d.dbNameCol())
}

func (Mysql) tableNameTpl() string {
	return common.ReduceSpaces(`
			SELECT table_name, table_schema, table_rows, table_comment
			FROM
				information_schema.tables
			WHERE
				table_type = '%s' AND table_schema %s
		`)
}

func (d Mysql) TableNameSql(dbname string, more bool) string {
	if more {
		dbcond := "LIKE " + common.WrapWith(dbname, "'", "%'")
		return fmt.Sprintf(d.tableNameTpl(), "BASE TABLE", dbcond)
	} else {
		dbcond := "= " + d.dbNameVal(dbname)
		return fmt.Sprintf(d.tableNameTpl(), "BASE TABLE", dbcond)
	}
}

func (d Mysql) ViewNameSql(dbname string, more bool) string {
	if more {
		dbcond := "LIKE " + common.WrapWith(dbname, "'", "%'")
		return fmt.Sprintf(d.tableNameTpl(), "VIEW", dbcond)
	} else {
		dbcond := "= " + d.dbNameVal(dbname)
		return fmt.Sprintf(d.tableNameTpl(), "VIEW", dbcond)
	}
}

func (Mysql) ColumnTypeSql(fullTableName string) string {
	return fmt.Sprintf("SELECT * FROM %s LIMIT 0", fullTableName)
}

func (d Mysql) ColumnInfoSql(table, dbname string) string {
	tpl := common.ReduceSpaces(`
			SELECT 
				column_name, column_type, column_key, 
				column_default, extra, column_comment,
				character_maximum_length,
				numeric_precision, numeric_scale
			FROM
				information_schema.columns
			WHERE
				table_name = '%s' AND table_schema = %s
			ORDER BY ordinal_position
		`)
	dbname = d.dbNameVal(dbname)
	return fmt.Sprintf(tpl, table, dbname)
}
