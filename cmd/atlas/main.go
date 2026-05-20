package main

import (
	"ariga.io/atlas-provider-gorm/gormschema"
	"io"
	"os"
	"todo/model"
)

func main() {
	stmts, err := gormschema.New("postgres").Load(
		&model.Todo{},
		// 追加モデルをここに列挙
	)
	if err != nil {
		os.Stderr.WriteString(err.Error())
		os.Exit(1)
	}
	io.WriteString(os.Stdout, stmts)
}
