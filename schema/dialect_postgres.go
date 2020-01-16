package schema

import (
	"fmt"

	"github.com/azhai/gozzo-utils/common"
	"github.com/azhai/gozzo-utils/redisw"
)

const PGSQL_DEFAULT_PORT uint16 = 5432

type Postgres struct {
}

func (Postgres) GetDSN(params redisw.ConnParams) (string, string) {
	dsn := "user=" + params.Username
	if params.Password != "" {
		dsn += " password=" + params.Password
	}
	if params.Host != "" {
		dsn += " host=" + params.Host
	}
	if port := params.StrPort(PGSQL_DEFAULT_PORT); port != "" {
		dsn += " port=" + port
	}
	if params.Database != "" {
		dsn += " dbname=" + params.Database
	}
	return "postgres", dsn
}

func (Postgres) QuoteIdent(ident string) string {
	return common.WrapWith(ident, `"`, `"`)
}

func (Postgres) dbNameCol() string {
	return "CURRENT_SCHEMA()"
}

func (d Postgres) dbNameVal(dbname string) string {
	if dbname == "" {
		return d.dbNameCol()
	} else {
		return common.WrapWith(dbname, "'", "'")
	}
}

func (d Postgres) CurrDbNameSql() string {
	return fmt.Sprintf("SELECT %s", d.dbNameCol())
}

func (Postgres) tableNameTpl() string {
	return common.ReduceSpaces(`
			SELECT table_name, table_schema, table_rows, table_comment
			FROM
				information_schema.tables
			WHERE
				table_type = '%s' AND table_schema = %s
		`)
}

func (d Postgres) TableNameSql(dbname string, more bool) string {
	if more {
		dbcond := "LIKE " + common.WrapWith(dbname, "'", "%'")
		return fmt.Sprintf(d.tableNameTpl(), "BASE TABLE", dbcond)
	} else {
		dbcond := "= " + d.dbNameVal(dbname)
		return fmt.Sprintf(d.tableNameTpl(), "BASE TABLE", dbcond)
	}
}

func (d Postgres) ViewNameSql(dbname string, more bool) string {
	if more {
		dbcond := "LIKE " + common.WrapWith(dbname, "'", "%'")
		return fmt.Sprintf(d.tableNameTpl(), "VIEW", dbcond)
	} else {
		dbcond := "= " + d.dbNameVal(dbname)
		return fmt.Sprintf(d.tableNameTpl(), "VIEW", dbcond)
	}
}

func (Postgres) ColumnTypeSql(fullTableName string) string {
	return fmt.Sprintf("SELECT * FROM %s LIMIT 0", fullTableName)
}

func (d Postgres) ColumnInfoSql(table, dbname string) string {
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
