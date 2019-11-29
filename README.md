# gozzo 尜舟

为gorm添加一些扩展功能，主要是从数据表生成 Model，或者反向操作，包括代码注释中获取字段说明。

## Windows 下编译

双击执行 WinBuild.bat 编译生成 dbtool.exe

## Linux/MacOS 下编译

在目录下使用 make 命令编译生成 dbtool

## 执行命令

除了可运行的程序 dbtool 或 dbtool.exe ，还要以下文件：
* 配置文件，默认 settings.toml
* 模板文件 gen_init.tmpl 和 gen_model.tmpl

运行命令从数据表中生成 models
```bash
dbtool -f settings.toml -d default
```
