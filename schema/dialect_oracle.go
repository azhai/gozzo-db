package schema

import (
	"fmt"

	"github.com/azhai/gozzo-utils/common"
	"github.com/azhai/gozzo-utils/redisw"
)

const ORACLE_DEFAULT_PORT uint16 = 1521

type Oracle struct {
}

func (Oracle) GetDSN(params redisw.ConnParams) (string, string) {
	user := redisw.ConcatWith(params.Username, params.Password)
	dsn := "oracle://"
	if user != "" {
		dsn += user + "@"
	}
	dsn += params.GetAddr("127.0.0.1", ORACLE_DEFAULT_PORT)
	if params.Database != "" {
		dsn += "?database" + params.Database
	}
	return "oracle", dsn
}

func (Oracle) QuoteIdent(ident string) string {
	return common.WrapWith(ident, "{", "}")
}

func (Oracle) dbNameVal(dbname string) string {
	if dbname == "" {
		return "SELECT SYS_CONTEXT('userenv', 'current_schema') FROM DUAL"
	} else {
		return common.WrapWith(dbname, "'", "'")
	}
}

func (Oracle) CurrDbNameSql() string {
	return "SELECT SYS_CONTEXT('db_name') FROM DUAL"
}

func (Oracle) tableNameTpl() string {
	return "SELECT %s, owner, num_rows* FROM %s WHERE owner %s"
}

func (d Oracle) TableNameSql(dbname string, more bool) string {
	if more {
		dbcond := "LIKE " + common.WrapWith(dbname, "'", "%'")
		return fmt.Sprintf(d.tableNameTpl(), "table_name", "all_tables", dbcond)
	} else {
		dbcond := "IN (" + d.dbNameVal(dbname) + ")"
		return fmt.Sprintf(d.tableNameTpl(), "table_name", "all_tables", dbcond)
	}
}

func (d Oracle) ViewNameSql(dbname string, more bool) string {
	if more {
		dbcond := "LIKE " + common.WrapWith(dbname, "'", "%'")
		return fmt.Sprintf(d.tableNameTpl(), "view_name", "all_views", dbcond)
	} else {
		dbcond := "IN (" + d.dbNameVal(dbname) + ")"
		return fmt.Sprintf(d.tableNameTpl(), "view_name", "all_views", dbcond)
	}
}

func (Oracle) ColumnTypeSql(fullTableName string) string {
	return fmt.Sprintf("SELECT * FROM %s WHERE 1=0", fullTableName)
}

func (d Oracle) ColumnInfoSql(table, dbname string) string {
	tpl := common.ReduceSpaces(`
			SELECT 
				column_name, column_type, column_key,
				column_default, extra, column_comment
			FROM
				cols
			WHERE
				table_name = '%s' AND table_schema = %s
			ORDER BY ordinal_position
		`)
	dbname = d.dbNameVal(dbname)
	return fmt.Sprintf(tpl, table, dbname)
}
