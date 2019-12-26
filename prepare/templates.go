package prepare

var templates = map[string]string{
	// table 文件模板
	"gen_table.tmpl": `
{{$rules := .Rules}}
{{$np := .NullPointer}}
// {{.Table.TableComment}}
type {{.Name}} struct {
	BaseModel
{{- range .Columns }}
	{{- if eq .FieldName "id" }}
	{{- else }}
		{{- $rule := GetRule $rules .FieldName}}
		{{GenNameType . $rule $np}} {{GenTagComment . $rule}}
	{{- end }}
{{- end }}
}

// 数据表名为 {{.Table.TableName}}
func ({{.Name}}) TableName() string {
	return "{{.Table.TableName}}"
}

// 数据表备注
func ({{.Name}}) TableComment() string {
	return "{{.Table.TableComment}}"
}`,


	// query 文件模板
	"gen_query.tmpl": `
// 查询符合条件的所有行
func (m {{.Name}}) FindAll(filters ...base.FilterFunc) (objs []*{{.Name}}, err error) {
	err = db.Model(m).Scopes(filters...).Find(&objs).Error
	err = IgnoreNotFoundError(err)
	return
}

// 查询符合条件的第一行
func (m {{.Name}}) GetOne(filters ...base.FilterFunc) (obj *{{.Name}}, err error) {
	obj = new({{.Name}})
	err = db.Model(m).Scopes(filters...).Take(&obj).Error
	err = IgnoreNotFoundError(err)
	return
}`,


	// init 文件模板
	"gen_init.tmpl": `
var (
	db *gorm.DB // 数据库对象
	ModelInsts = []interface{}{ // 所有Model实例
		{{.Models}}
	}
)

type BaseModel = base.Model

// 忽略表中无数据的错误
func IgnoreNotFoundError(err error) error {
	return base.IgnoreNotFoundError(err)
}

// 获取当前db
func Query() *gorm.DB {
	return db
}

// 查询某张数据表
func QueryTable(name string) *gorm.DB {
	return db.Table(name)
}

// 连接数据库
func init() {
	conf, err := prepare.GetConfig("{{.FileName}}")
	if err != nil {
		panic(err)
	}
	if c, ok := conf.Connections["cache"]; ok && c.Driver == "redis" {
		rds := cache.ConnectRedisPool(c.ConnParams)
		cache.SetRedisBackend(rds)
	}
	db, err = gorm.Open(conf.GetDSN("{{.ConnName}}"))
	if err != nil {
		panic(err)
	}

	// 初始化数据库
	if conf.Application.Debug {
		db = db.Debug().LogMode(true)
		db.SetLogger(log.New(os.Stdout, "\r\n", 0))
	}
	drv := conf.GetDriverName("{{.ConnName}}")
	if drv == "mysql" {
		db.Set("gorm:table_options", "ENGINE=InnoDB")
	}
	db = MigrateTables(drv, db)
}

// 自动建表，如果缺少表或字段会加上
func MigrateTables(drv string, db *gorm.DB) *gorm.DB {
	{{- if .Plural }}
		db.SingularTable(false)
	{{- else }}
		db.SingularTable(true)
	{{- end }}
	db = db.AutoMigrate(ModelInsts...) // 创建缺少的表和字段
	if drv == "mysql" { // 更新MySQL表注释
		db = export.AlterTableComments(db, ModelInsts...)
	}
	return db
}`,
}
