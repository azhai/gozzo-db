package main

import (
	"flag"
	"fmt"

	"github.com/azhai/gozzo-db/cmd"
	"github.com/azhai/gozzo-db/construct"
)

var (
	/**
	 * 优化Model结构，可选 0-2
	 * 0  不优化，只输出信息
	 * 1  一般优化
	 * 2  尽力优化
	 */
	level   uint
	verbose bool // 详细输出
)

// 优化Model结构
func main() {
	flag.UintVar(&level, "level", 0, "优化等级，可选 0-2 ，具体参考说明")
	flag.BoolVar(&verbose, "vv", false, "输出详细信息")
	conf := cmd.Initialize(func() error {
		if level >= 3 {
			return fmt.Errorf("没有这种等级，可选 0-2 ，具体参考说明")
		}
		return nil
	})

	outDir := conf.Application.OutputDir
	construct.ScanModelDir(outDir)
	if verbose {
		for _, c := range construct.ModelClasses {
			fmt.Println(c.Name, ":")
			fmt.Println(c.GetInnerCode())
		}
	}
}
