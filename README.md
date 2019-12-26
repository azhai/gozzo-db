# gozzo 尜舟

为 gorm 添加一些扩展功能，主要是从数据表生成 Model，或者反向操作，包括代码注释中获取字段说明。

## 编译

Windows 下编译：

双击执行 WinBuild.bat

编译生成 gen2model.exe sync2table.exe dump2file.exe

Linux/MacOS 下编译：

在目录下使用 make 命令

编译生成 gen2model sync2table dump2file

## gen2model

根据数据表结构生成对应的 Model 代码

除了可运行的程序 gen2model 或 gen2model.exe ，还要以下文件：
* 配置文件，默认 settings.toml
* 模板文件（可选） gen_init.tmpl gen_table.tmpl gen_query.tmpl

运行命令从数据表中生成 models

```bash
gen2model -f settings.toml -d default -mode 0 -v
```

生成文件不同结构，可选 0-5
* 0  与 5 类似，但会在每个文件名前面加一个下划线
* 1  只生成 init.go 文件
* 2  除了 init.go 文件， table 和 query 都放入 tables.go 中
* 3  除了 init.go 文件， table 都放入 tables.go 中， query 都放入 queries.go 中
* 4  除了 init.go 文件， table 都放入 tables.go 中， query 分开放入对应模型名文件中
* 5  除了 init.go 文件， table 和 query 一起放入对应模型名文件中

## sync2table

从代码中同步到数据表中，包括缺少的字段、索引和改动的注释

NOTE: 编译依赖于 gen2model 生成的 models

## dump2file

在数据表和TOML文件之间导入导出数据

NOTE: 编译依赖于 gen2model 生成的 models

## 初次使用

先修改 settings.toml 中的数据库连接参数，由于后两个命令依赖于 gen2model 生成的代码，建议如下操作

```bash
make  #第一遍编译，尚未有任何 Model 实例，只有 gen2model 命令可用
gen2model -f settings.toml -mode 3  #生成当前数据库的 Model 实例
make  #第二遍编译，sync2table 和 dump2file 命令可用了
#如果以后 models 中代码有变动，还需要重新编译
```
