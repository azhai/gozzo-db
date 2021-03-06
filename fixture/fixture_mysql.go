package fixture

import (
	"fmt"
	"testing"
	"time"

	"github.com/azhai/gozzo-db/prepare"
	"github.com/azhai/gozzo-utils/common"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

var (
	TestTableName      = "people"
	TestCreateTableSql = `CREATE TABLE IF NOT EXISTS [%s](
		[id] INT UNSIGNED NOT NULL AUTO_INCREMENT,
		[name] VARCHAR(128),
		[height] FLOAT, 
		[birth] DATETIME,
		PRIMARY KEY ([id]) USING BTREE
	)`
)

func InitDB() *gorm.DB {
	conf, _ := prepare.GetConfig("../settings.toml")
	db, err := gorm.Open(conf.GetDSN("default"))
	if err != nil {
		panic(err)
	}
	sql := common.ReplaceQuotes(TestCreateTableSql)
	db.Exec(fmt.Sprintf(sql, TestTableName+"_males"))
	db.Exec(fmt.Sprintf(sql, TestTableName+"_females"))
	if testing.Verbose() {
		return db.Debug() // 开启Debug模式输出SQL
	} else {
		return db
	}
}

func GetDate(dt string) (t time.Time) {
	t, _ = time.Parse("2006-01-02", dt)
	return
}

func TruncateRecords(db *gorm.DB) *gorm.DB {
	db.Exec(fmt.Sprintf("TRUNCATE TABLE %s", TestTableName+"_males"))
	db.Exec(fmt.Sprintf("TRUNCATE TABLE %s", TestTableName+"_females"))
	return db
}
