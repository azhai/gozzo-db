package prepare

import (
	"github.com/BurntSushi/toml"
	"github.com/azhai/gozzo-db/schema"
	"github.com/azhai/gozzo-db/utils"
)

// 解析配置和创建日志
func GetConfig(fileName string) (*Config, error) {
	var conf = new(Config)
	fullPath := utils.GetAbsFile(fileName)
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
	PluralTable bool   `toml:"plural_table"` // 表名使用复数形式
	NullPointer bool   `toml:"null_pointer"` // 字段可为NULL时，使用对应的指针类型
}

// 连接配置
type ConnConfig struct {
	Driver string `toml:"driver"`
	Prefix string `toml:"prefix"`
	schema.ConnParams
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
	FileName    string // 文件名
	ConnName    string // 连接名
	Application AppConfig
	Connections map[string]ConnConfig
	ModelRules  map[string]TableRuleConfig
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
