package prepare

import (
	"strings"

	"github.com/BurntSushi/toml"
	"github.com/azhai/gozzo-db/schema"
	"github.com/azhai/gozzo-utils/filesystem"
	"github.com/azhai/gozzo-utils/redisw"
)

// 解析配置和创建日志
func GetConfig(fileName string) (*Config, error) {
	var conf = new(Config)
	fullPath := filesystem.GetAbsFile(fileName)
	_, err := toml.DecodeFile(fullPath, &conf)
	if err != nil {
		return nil, err
	}
	conf.FileName = fileName
	if conf.ConnName == "" {
		conf.ConnName = "default"
	}
	return conf, nil
}

/**
***********************************************************
* 配置解析
***********************************************************
**/

// 应用配置
type AppConfig struct {
	Debug       bool   `toml:"debug"`
	OutputDir   string `toml:"output_dir"`   // 输出目录，例如 models
	DataFile    string `toml:"data_file"`    // 数据文件，例如 data.toml
	PluralTable bool   `toml:"plural_table"` // 表名使用复数形式
}

// 连接配置
type ConnConfig struct {
	Driver string `toml:"driver"`
	Prefix string `toml:"prefix"`
	redisw.ConnParams
}

// NULL字段使用指针类型
type NullPointer struct {
	UsePointer bool `toml:"use_pointer"` // 是否使用指针类型
	MustIndex  bool `toml:"must_index"`  // 如果使用指针类型，字段必须是索引
	MinLength  int  `toml:"min_length"`  // 如果使用指针类型，字段长度最少为多大
}

func (np NullPointer) MatchCond(ci *schema.ColumnInfo) bool {
	if np.MustIndex && !ci.IsIndex() {
		return false
	}
	dbtype := strings.ToLower(ci.DatabaseTypeName())
	if strings.HasSuffix(dbtype, "text") {
		return true // TEXT字段
	} else if strings.HasSuffix(dbtype, "varchar") {
		var size int
		if size = ci.GetSize(); size == 0 {
			size = 255
		}
		return size < 0 || size >= np.MinLength
	}
	return false
}

func NullPointerMatch(nps map[string]NullPointer, rule RuleConfig, ci *schema.ColumnInfo) bool {
	switch rule.Type {
	case "int", "uint", "int64", "uint64":
		if np, ok := nps["int"]; ok && np.UsePointer {
			return np.MatchCond(ci)
		}
	case "string":
		if np, ok := nps["string"]; ok && np.UsePointer {
			return np.MatchCond(ci)
		}
	case "time.Time":
		if np, ok := nps["time"]; ok && np.UsePointer {
			return np.MatchCond(ci)
		}
	}
	return false
}

// 规则配置
type RuleConfig struct {
	Name    string `toml:"name"`
	Type    string `toml:"type"`
	Json    string `toml:"json"`
	Tags    string `toml:"tags"`
	Comment string `toml:"comment"`
}

type TableRuleConfig = map[string]RuleConfig

func GetRule(rules TableRuleConfig, name string) RuleConfig {
	if rule, ok := rules[name]; ok {
		return rule
	}
	return RuleConfig{}
}

// 配置
type Config struct {
	FileName     string // 文件名
	ConnName     string // 连接名
	Application  AppConfig
	Connections  map[string]ConnConfig
	NullPointers map[string]NullPointer     `toml:"null_pointers"`
	ModelRules   map[string]TableRuleConfig `toml:"model_rules"`
}

func (c Config) GetDriverName(name string) string {
	if params, ok := c.Connections[name]; ok {
		if params.Driver == "sqlite3" {
			return "sqlite" // Sqlite的import包名和Open()驱动名不一样
		}
		return params.Driver
	}
	return ""
}

func (c Config) GetTablePrefix(name string) string {
	if params, ok := c.Connections[name]; ok {
		return params.Prefix
	}
	return ""
}

func (c Config) GetDSN(name string) (string, string) {
	if params, ok := c.Connections[name]; ok {
		dia := schema.GetDialectByName(params.Driver)
		if dia != nil {
			return dia.GetDSN(params.ConnParams)
		}
	}
	return name, ""
}

func (c Config) GetRules(table string) (rules TableRuleConfig) {
	if baseRules, ok := c.ModelRules["_"]; ok {
		rules = baseRules
	} else {
		rules = TableRuleConfig{}
	}
	if tableRules, ok := c.ModelRules[table]; ok {
		for name, colRule := range tableRules {
			rules[name] = colRule
		}
	}
	return
}
