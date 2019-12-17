package schema

import (
	"fmt"

	"github.com/azhai/gozzo-db/utils"
)

type Sqlite struct {
}

func (Sqlite) GetDSN(params ConnParams) (string, string) {
	user := params.Concat(params.Username, params.Password)
	var dsn string
	if user != "" {
		dsn = user + "@"
	}
	dsn += params.Database + "?cache=shared&mode=rwc"
	return "sqlite3", dsn
}

func (Sqlite) QuoteIdent(ident string) string {
	return utils.WrapWith(ident, "`", "`")
}

func (Sqlite) dbNameVal(dbname string) string {
	if dbname == "" {
		return "'main'"
	} else {
		return utils.WrapWith(dbname, "'", "'")
	}
}

func (Sqlite) CurrDbNameSql() string {
	return "PRAGMA database_list"
}

func (Sqlite) tableNameTpl() string {
	return utils.ReduceSpaces(`
			SELECT name
			FROM
				sqlite_master
			WHERE
				type ='%s' AND name NOT LIKE 'sqlite_%%'
		`)
}

func (d Sqlite) TableNameSql(dbname string, more bool) string {
	return fmt.Sprintf(d.tableNameTpl(), "table")
}

func (d Sqlite) ViewNameSql(dbname string, more bool) string {
	return fmt.Sprintf(d.tableNameTpl(), "view")
}

func (Sqlite) ColumnTypeSql(fullTableName string) string {
	return fmt.Sprintf("SELECT * FROM %s LIMIT 0", fullTableName)
}

func (d Sqlite) ColumnInfoSql(table, dbname string) string {
	tpl := "PRAGMA table_info(%s)"
	return fmt.Sprintf(tpl, table)
}
