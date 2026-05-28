package main

import (
	"io"
	"os"
	"todo/model"

	"ariga.io/atlas-provider-gorm/gormschema"
)

func main() {
	stmts, err := gormschema.New("postgres").Load(
		&model.Company{},
		&model.User{},
		&model.CompanyMember{},
		&model.Team{},
		&model.TeamMember{},
		&model.Todo{},
		// 追加モデルをここに列挙
	)
	if err != nil {
		os.Stderr.WriteString(err.Error())
		os.Exit(1)
	}
	io.WriteString(os.Stdout, stmts)
}
