package export

import (
	"bytes"
	"fmt"
	"io/ioutil"

	"github.com/BurntSushi/toml"
	"github.com/azhai/gozzo-db/construct"
	"github.com/azhai/gozzo-db/schema"
	"github.com/jinzhu/gorm"
)

var tomlBuffer = new(bytes.Buffer)

type Dict = map[string]interface{}

// 将Dict转为具体对象，键名与属性名匹配时，使用toml标签或不区分大小写的方式
func ConvertToml(src Dict, dst interface{}) error {
	tomlBuffer.Reset()
	err := toml.NewEncoder(tomlBuffer).Encode(src)
	if err != nil {
		return err
	}
	_, err = toml.DecodeReader(tomlBuffer, dst)
	return err
}

// 往数据库中写入数据
func InsertRows(scope *gorm.Scope, rows []Dict, verbose bool) error {
	db := scope.DB()
	for _, row := range rows {
		if err := ConvertToml(row, scope.Value); err != nil {
			return err
		}
		if field := scope.PrimaryField(); field != nil {
			field.Set(0) // 清空主键，避免下面Create变成更新
		}
		if verbose {
			fmt.Printf("%#v\n", scope.Value)
		}
		db = db.Create(scope.Value)
	}
	return db.Error
}

// 将toml文件中数据导入数据库，使用表名作为section
func LoadFileData(db *gorm.DB, fname string, models []interface{}, verbose bool) (*gorm.DB, error) {
	temp := make(map[string][]Dict)
	if _, err := toml.DecodeFile(fname, &temp); err != nil {
		return db, err
	}
	sch := schema.NewSchema(db.DB())
	tbInfos := sch.ListTable("", false)
	for _, m := range models {
		scope := db.NewScope(m)
		tableName := scope.TableName()
		if rows, ok := temp[tableName]; ok {
			info, ok := tbInfos[tableName]
			if !ok || info.TableRows == int64(len(rows)) {
				continue // 可能会重复导入
			}
			if verbose {
				fmt.Printf("Insert %d rows into table %s\n", len(rows), tableName)
			}
			_ = InsertRows(scope, rows, verbose)
		}
	}
	return db, nil
}

func GetTableName(obj interface{}, name string) string {
	if m, ok := obj.(construct.ITableName); ok {
		return m.TableName()
	}
	return name
}

type Exportor struct {
	Data map[string][]interface{}
}

func NewExportor() *Exportor {
	return &Exportor{
		Data: make(map[string][]interface{}),
	}
}

func (ep *Exportor) AddObject(obj interface{}, group string) bool {
	if obj == nil {
		return false
	}
	if group = GetTableName(obj, group); group == "" {
		return false
	}
	ep.Data[group] = append(ep.Data[group], obj)
	return true
}

func (ep *Exportor) SetBuffer(buf *bytes.Buffer) error {
	if len(ep.Data) > 0 {
		return toml.NewEncoder(buf).Encode(ep.Data)
	}
	return nil
}

func (ep *Exportor) WriteTo(fname string) error {
	buf := new(bytes.Buffer)
	if err := ep.SetBuffer(buf); err != nil {
		return err
	}
	return ioutil.WriteFile(fname, buf.Bytes(), 0644)
}
