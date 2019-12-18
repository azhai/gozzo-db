package main

import (
	"fmt"

	"github.com/azhai/gozzo-db/cmd"
	"github.com/azhai/gozzo-db/prepare"
)

func main() {
	conf, db := cmd.Initialize()
	names, err := prepare.CreateModels(conf, db)
	if err != nil {
		fmt.Println(err)
	}
	err = prepare.GenInitFile(conf, names)
	if err != nil {
		fmt.Println(err)
	}
}
