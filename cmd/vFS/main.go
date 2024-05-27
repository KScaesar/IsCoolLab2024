package main

import (
	"github.com/KScaesar/IsCoolLab2024/pkg/adapters/database"
	"github.com/KScaesar/IsCoolLab2024/pkg/inject"
)

func main() {
	conf := &database.GormConfing{
		Dsn:     "vFS.db",
		Migrate: true,
	}

	infra, err := inject.NewInfra(conf)
	if err != nil {
		panic(err)
	}

	command := inject.NewRootCommand(infra)

	command.Execute()

	infra.Cleanup()
}
