package schema

import (
	"fmt"

	"github.com/azhai/gozzo-db/utils"
)

type Postgres struct {
}

func (Postgres) GetDSN(params ConnParams) (string, string) {
	dsn := "user=" + params.Username
	if params.Password != "" {
		dsn += " password=" + params.Password
	}
	if params.Host != "" {
		dsn += " host=" + params.Host
	}
	if port := params.StrPort(); port != "" {
		dsn += " port=" + port
	}
	if params.Database != "" {
		dsn += " dbname=" + params.Database
	}
	return "postgres", dsn
}

func (Postgres) QuoteIdent(ident string) string {
	return utils.WrapWith(ident, `"`, `"`)
}

func (Postgres) dbNameCol() string {
	return "CURRENT_SCHEMA()"
}

func (d Postgres) dbNameVal(dbname string) string {
	if dbname == "" {
		return d.dbNameCol()
	} else {
		return utils.WrapWith(dbname, "'", "'")
	}
}

func (d Postgres) CurrDbNameSql() string {
	return fmt.Sprintf("SELECT %s", d.dbNameCol())
}

func (Postgres) tableNameTpl() string {
	return utils.TrimTail(`
			SELECT table_name, table_comment, table_rows
			FROM
				information_schema.tables
			WHERE
				table_type = '%s' AND table_schema = %s
		`)
}

func (d Postgres) TableNameSql(dbname string) string {
	dbname = d.dbNameVal(dbname)
	return fmt.Sprintf(d.tableNameTpl(), "BASE TABLE", dbname)
}

func (d Postgres) ViewNameSql(dbname string) string {
	dbname = d.dbNameVal(dbname)
	return fmt.Sprintf(d.tableNameTpl(), "VIEW", dbname)
}

func (Postgres) ColumnTypeSql(fullTableName string) string {
	return fmt.Sprintf("SELECT * FROM %s LIMIT 0", fullTableName)
}

func (d Postgres) ColumnInfoSql(table, dbname string) string {
	tpl := utils.TrimTail(`
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
