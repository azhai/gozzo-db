package schema

import (
	"fmt"

	"github.com/azhai/gozzo-db/utils"
)

type Mssql struct {
}

func (Mssql) GetDSN(params ConnParams) (string, string) {
	user := params.Concat(params.Username, params.Password)
	addr := params.Concat(params.Host, params.StrPort())
	dsn := "sqlserver://"
	if user != "" {
		dsn += user + "@"
	}
	dsn += addr
	if params.Database != "" {
		dsn += "?database" + params.Database
	}
	return "mssql", dsn
}

func (Mssql) QuoteIdent(ident string) string {
	return utils.WrapWith(ident, "[", "]")
}

func (Mssql) dbNameCol() string {
	return "DB_NAME()"
}

func (d Mssql) dbNameVal(dbname string) string {
	if dbname == "" {
		return d.dbNameCol()
	} else {
		return utils.WrapWith(dbname, "'", "'")
	}
}

func (d Mssql) CurrDbNameSql() string {
	return fmt.Sprintf("SELECT %s", d.dbNameCol())
}

func (Mssql) tableNameTpl() string {
	return utils.TrimTail(`
			SELECT T.name as name
			FROM
				sys.%s AS T
				INNER JOIN sys.schemas AS S ON S.schema_id = T.schema_id
				LEFT JOIN sys.extended_properties AS EP ON EP.major_id = T.[object_id]
			WHERE
				T.is_ms_shipped = 0 AND
				(EP.class_desc IS NULL OR (EP.class_desc <> 'OBJECT_OR_COLUMN' AND
				EP.[name] <> 'microsoft_database_tools_support'))
		`)
}

func (d Mssql) TableNameSql(dbname string) string {
	return fmt.Sprintf(d.tableNameTpl(), "tables")
}

func (d Mssql) ViewNameSql(dbname string) string {
	return fmt.Sprintf(d.tableNameTpl(), "views")
}

func (Mssql) ColumnTypeSql(fullTableName string) string {
	return fmt.Sprintf("SELECT * FROM %s WHERE 1=0", fullTableName)
}

func (d Mssql) ColumnInfoSql(table, dbname string) string {
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
