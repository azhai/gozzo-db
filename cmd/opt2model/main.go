package main

import (
	"flag"
	"fmt"

	"github.com/azhai/gozzo-db/cmd"
	"github.com/azhai/gozzo-db/construct"
)

/**
 * 优化Model结构，可选 0-2
 * 0  不优化，只输出信息
 * 1  一般优化
 * 2  尽力优化
 */
var level uint

// 优化Model结构
func main() {
	flag.UintVar(&level, "level", 0, "优化等级，可选 0-2 ，具体参考说明")
	conf := cmd.Initialize(func() error {
		if level >= 3 {
			return fmt.Errorf("没有这种等级，可选 0-2 ，具体参考说明")
		}
		return nil
	})

	outDir := conf.Application.OutputDir
	construct.ScanModelDir(outDir)
	for _, c := range construct.ModelClasses {
		fmt.Println(c.Name, ":")
		fmt.Println(c.GetInnerCode())
	}
}
