package prepare

import (
	"github.com/BurntSushi/toml"
	"github.com/azhai/gozzo-db/schema"
)

// 解析配置和创建日志
func GetConfig(filename string) *Config {
	var conf = &Config{}
	_, err := toml.DecodeFile(filename, &conf)
	if err == nil {
		return conf
	}
	return nil
}

/**
***********************************************************
* 配置解析
***********************************************************
**/

// 应用配置
type AppConfig struct {
	OutputDir   string `toml:"output_dir"`
	TablePrefix string `toml:"table_prefix"`
	PluralTable bool   `toml:"plural_table"`
}

// 连接配置
type ConnConfig struct {
	Driver string `toml:"driver"`
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

// 配置
type Config struct {
	Application AppConfig
	Connections map[string]ConnConfig
	ModelRules  map[string]TableRuleConfig
}

func (c Config) GetDSN(name string) (string, string) {
	if params, ok := c.Connections[name]; ok {
		name = params.Driver
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
